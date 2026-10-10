package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sujanv/LogSift/pkg/filter"
	"github.com/sujanv/LogSift/pkg/model"
	"github.com/sujanv/LogSift/pkg/output"
	"github.com/sujanv/LogSift/pkg/parser"
	"github.com/sujanv/LogSift/pkg/reader"
	"github.com/sujanv/LogSift/pkg/stats"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcmd := os.Args[1]
	switch subcmd {
	case "search":
		runSearch(os.Args[2:])
	case "summarize", "stats":
		runSummarize(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		// Default to search if first arg is not a command
		runSearch(os.Args[1:])
	}
}

func printUsage() {
	fmt.Println(`LogSift - High-throughput structured log analysis tool

Usage:
  logsift <command> [options] [file...]

Commands:
  search      Search and filter log streams in real-time
  summarize   Generate summary metrics (latency, error clusters, HTTP status)
  help        Show help for LogSift

Examples:
  logsift search --level ERROR /var/log/app.log
  logsift search --filter "status>=500 AND latency>200ms" app.log.gz
  cat access.log | logsift summarize
  logsift summarize --report report.md app.log`)
}

func runSearch(args []string) {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	levelFilter := fs.String("level", "", "Filter by log level (DEBUG, INFO, WARN, ERROR, FATAL)")
	exprFilter := fs.String("filter", "", "Filter expression e.g. 'status>=500 AND latency>200ms'")
	sinceFilter := fs.String("since", "", "Filter logs within relative time duration (e.g. 15m, 2h, 1d)")
	format := fs.String("format", "auto", "Log format: auto, json, logfmt, clf, syslog")
	outMode := fs.String("output", "terminal", "Output mode: terminal, json, csv")
	noColor := fs.Bool("no-color", false, "Disable colored ANSI output")

	_ = fs.Parse(args)
	files := fs.Args()

	var pred filter.Predicate
	if *exprFilter != "" {
		p, err := filter.ParseExpression(*exprFilter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid filter expression: %v\n", err)
			os.Exit(1)
		}
		pred = p
	}

	var timeMatcher *filter.TimeRangeMatcher
	if *sinceFilter != "" {
		from, err := filter.ParseSince(time.Now(), *sinceFilter)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid --since duration: %v\n", err)
			os.Exit(1)
		}
		timeMatcher = filter.NewTimeRangeMatcher(from, time.Time{})
	}

	r := getReader(files)
	defer r.Close()

	sr := reader.NewStreamReader(r)
	defer sr.Close()

	logParser := selectParser(*format, sr)
	pool := model.DefaultRecordPool

	termOut := output.NewTerminalFormatter(os.Stdout, *noColor)
	jsonOut := output.NewJSONExporter(os.Stdout)
	csvOut := output.NewCSVExporter(os.Stdout)
	defer csvOut.Flush()

	targetLevel := model.ParseLevel(*levelFilter)

	for {
		line, lineNum, offset, err := sr.NextLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Fprintf(os.Stderr, "read error at line %d: %v\n", lineNum, err)
			break
		}

		if len(strings.TrimSpace(string(line))) == 0 {
			continue
		}

		rec := pool.Get()
		rec.LineNum = lineNum
		rec.Offset = offset

		if err := logParser.Parse(line, rec); err != nil {
			rec.Message = string(line)
		}

		// Apply filters
		if targetLevel != model.LevelUnknown && rec.Level < targetLevel {
			pool.Put(rec)
			continue
		}

		if timeMatcher != nil && !timeMatcher.Match(rec) {
			pool.Put(rec)
			continue
		}

		if pred != nil && !pred.Match(rec) {
			pool.Put(rec)
			continue
		}

		switch *outMode {
		case "json":
			_ = jsonOut.Export(rec)
		case "csv":
			_ = csvOut.Export(rec)
		default:
			termOut.FormatRecord(rec)
		}

		pool.Put(rec)
	}
}

func runSummarize(args []string) {
	fs := flag.NewFlagSet("summarize", flag.ExitOnError)
	reportFile := fs.String("report", "", "Write markdown report to specified file")
	format := fs.String("format", "auto", "Log format: auto, json, logfmt, clf, syslog")
	noColor := fs.Bool("no-color", false, "Disable colored ANSI output")

	_ = fs.Parse(args)
	files := fs.Args()

	r := getReader(files)
	defer r.Close()

	sr := reader.NewStreamReader(r)
	defer sr.Close()

	logParser := selectParser(*format, sr)
	pool := model.DefaultRecordPool

	latTracker := stats.NewLatencyTracker()
	errSummarizer := stats.NewErrorSummarizer()
	httpTracker := stats.NewHTTPAnalytics()
	histogram := stats.NewTimeHistogram(time.Minute)

	startTime := time.Now()
	var totalRecords int64
	var totalErrors int64

	for {
		line, lineNum, offset, err := sr.NextLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			break
		}

		if len(strings.TrimSpace(string(line))) == 0 {
			continue
		}

		totalRecords++
		rec := pool.Get()
		rec.LineNum = lineNum
		rec.Offset = offset

		_ = logParser.Parse(line, rec)

		isErr := rec.Level >= model.LevelError || rec.StatusCode >= 500
		if isErr {
			totalErrors++
			errSummarizer.Add(rec.Message, rec.Timestamp)
		}

		if rec.Latency > 0 {
			latTracker.Add(rec.Latency)
		}

		if rec.StatusCode > 0 || rec.Path != "" {
			httpTracker.Add(rec.StatusCode, rec.Path, rec.Latency)
		}

		if !rec.Timestamp.IsZero() {
			histogram.Add(rec.Timestamp, isErr)
		}

		pool.Put(rec)
	}

	elapsed := time.Since(startTime)
	latSnap := latTracker.Snapshot()
	topErrors := errSummarizer.TopClusters(5)
	statusBD := httpTracker.StatusBreakdown()
	slowEP := httpTracker.TopSlowestEndpoints(5)

	if *reportFile != "" {
		f, err := os.Create(*reportFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create report file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		mdGen := output.NewMarkdownReportGenerator(f)
		_ = mdGen.Generate(totalRecords, elapsed, latSnap, topErrors, statusBD, slowEP)
		fmt.Printf("Summary report generated: %s\n", *reportFile)
		return
	}

	// Print terminal summary
	term := output.NewTerminalFormatter(os.Stdout, *noColor)
	fmt.Printf("\n=== LogSift Summary ===\n")
	fmt.Printf("Records: %d | Errors: %d | Processed in: %v\n\n", totalRecords, totalErrors, elapsed)

	fmt.Println("Latency Metrics:")
	latHeaders := []string{"Min", "Mean", "P50", "P75", "P90", "P95", "P99", "Max"}
	latRow := []string{
		latSnap.Min.String(),
		latSnap.Mean.String(),
		latSnap.P50.String(),
		latSnap.P75.String(),
		latSnap.P90.String(),
		latSnap.P95.String(),
		latSnap.P99.String(),
		latSnap.Max.String(),
	}
	term.PrintTable(latHeaders, [][]string{latRow})

	if len(topErrors) > 0 {
		fmt.Printf("\nTop Errors:\n")
		errHeaders := []string{"Count", "Sample Error"}
		errRows := make([][]string, 0, len(topErrors))
		for _, e := range topErrors {
			errRows = append(errRows, []string{fmt.Sprintf("%d", e.Count), e.Sample})
		}
		term.PrintTable(errHeaders, errRows)
	}

	fmt.Printf("\nHTTP Status Breakdown: 2xx: %d | 3xx: %d | 4xx: %d | 5xx: %d\n",
		statusBD["2xx"], statusBD["3xx"], statusBD["4xx"], statusBD["5xx"])
}

func getReader(files []string) io.ReadCloser {
	if len(files) == 0 {
		return os.Stdin
	}

	f, err := os.Open(files[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open %s: %v\n", files[0], err)
		os.Exit(1)
	}
	return f
}

func selectParser(format string, sr *reader.StreamReader) parser.Parser {
	switch strings.ToLower(format) {
	case "json":
		return parser.NewJSONParser()
	case "logfmt":
		return parser.NewLogfmtParser()
	case "clf":
		return parser.NewCLFParser()
	case "syslog":
		return parser.NewSyslogParser()
	default:
		// Auto-detect format by inspecting initial line
		line, _, _, err := sr.NextLine()
		if err == nil && len(line) > 0 {
			p := parser.AutoDetect([][]byte{line})
			return p
		}
		return parser.NewJSONParser()
	}
}
