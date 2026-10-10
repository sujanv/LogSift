package main

import (
	"database/sql"
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
	"github.com/sujanv/LogSift/pkg/storage"
	"github.com/sujanv/LogSift/pkg/tui"
	_ "modernc.org/sqlite"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	subcmd := os.Args[1]
	switch subcmd {
	case "tui":
		runTUI(os.Args[2:])
	case "search":
		runSearch(os.Args[2:])
	case "summarize", "stats":
		runSummarize(os.Args[2:])
	case "index":
		runIndex(os.Args[2:])
	case "query", "sql":
		runQuery(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
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
  index       Index log file into SQLite for instant SQL querying
  query       Execute arbitrary SQL against indexed SQLite database
  help        Show help for LogSift

Examples:
  logsift search --level ERROR /var/log/app.log
  logsift search --filter "status>=500 AND latency>200ms" app.log.gz
  cat access.log | logsift summarize
  logsift index --db logs.db app.log
  logsift query --db logs.db "SELECT path, AVG(latency_ms), COUNT(*) FROM logs GROUP BY path ORDER BY 2 DESC LIMIT 5"`)
}

func runIndex(args []string) {
	fs := flag.NewFlagSet("index", flag.ExitOnError)
	dbFile := fs.String("db", "logs.db", "SQLite database destination")
	format := fs.String("format", "auto", "Log format")
	batchSize := fs.Int("batch", 10000, "Insert batch size")

	_ = fs.Parse(args)
	files := fs.Args()

	r := getReader(files)
	defer r.Close()

	sr := reader.NewStreamReader(r)
	defer sr.Close()

	logParser := selectParser(*format, sr)
	indexer, err := storage.OpenSQLite(*dbFile, *batchSize)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize sqlite indexer: %v\n", err)
		os.Exit(1)
	}
	defer indexer.Close()

	pool := model.DefaultRecordPool
	t0 := time.Now()
	var count int64

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

		count++
		rec := pool.Get()
		rec.LineNum = lineNum
		rec.Offset = offset

		_ = logParser.Parse(line, rec)
		if err := indexer.Index(rec); err != nil {
			fmt.Fprintf(os.Stderr, "index error at line %d: %v\n", lineNum, err)
			break
		}
		pool.Put(rec)
	}

	_ = indexer.Flush()
	elapsed := time.Since(t0)
	fmt.Printf("Indexed %d logs into %s in %v (%.0f logs/sec)\n", count, *dbFile, elapsed, float64(count)/elapsed.Seconds())
}

func runQuery(args []string) {
	fs := flag.NewFlagSet("query", flag.ExitOnError)
	dbFile := fs.String("db", "logs.db", "SQLite database source")
	_ = fs.Parse(args)

	queryArgs := fs.Args()
	if len(queryArgs) == 0 {
		fmt.Fprintln(os.Stderr, "error: missing SQL query string")
		os.Exit(1)
	}
	query := strings.Join(queryArgs, " ")

	db, err := sql.Open("sqlite", *dbFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open database %s: %v\n", *dbFile, err)
		os.Exit(1)
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "query error: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get columns: %v\n", err)
		os.Exit(1)
	}

	var tableRows [][]string
	for rows.Next() {
		rawVals := make([]interface{}, len(cols))
		scanArgs := make([]interface{}, len(cols))
		for i := range rawVals {
			scanArgs[i] = &rawVals[i]
		}
		if err := rows.Scan(scanArgs...); err != nil {
			continue
		}

		row := make([]string, len(cols))
		for i, v := range rawVals {
			if v == nil {
				row[i] = "NULL"
			} else {
				row[i] = fmt.Sprintf("%v", v)
			}
		}
		tableRows = append(tableRows, row)
	}

	term := output.NewTerminalFormatter(os.Stdout, false)
	term.PrintTable(cols, tableRows)
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
		line, _, _, err := sr.NextLine()
		if err == nil && len(line) > 0 {
			return parser.AutoDetect([][]byte{line})
		}
		return parser.NewJSONParser()
	}
}

func runTUI(args []string) {
	files := args
	r := getReader(files)
	defer r.Close()

	sr := reader.NewStreamReader(r)
	defer sr.Close()

	logParser := selectParser("auto", sr)
	var records []*model.Record
	count := 0

	for {
		line, lineNum, offset, err := sr.NextLine()
		if err != nil || count >= 5000 {
			break
		}
		if len(strings.TrimSpace(string(line))) == 0 {
			continue
		}

		rec := &model.Record{
			LineNum: lineNum,
			Offset:  offset,
			Raw:     make([]byte, len(line)),
		}
		copy(rec.Raw, line)
		_ = logParser.Parse(line, rec)
		records = append(records, rec)
		count++
	}

	_ = tui.Run(records)
}
