package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/sujanv/LogSift/pkg/correlate"
	"github.com/sujanv/LogSift/pkg/diff"
	"github.com/sujanv/LogSift/pkg/filter"
	"github.com/sujanv/LogSift/pkg/mask"
	"github.com/sujanv/LogSift/pkg/model"
	"github.com/sujanv/LogSift/pkg/output"
	"github.com/sujanv/LogSift/pkg/parser"
	"github.com/sujanv/LogSift/pkg/reader"
	"github.com/sujanv/LogSift/pkg/server"
	"github.com/sujanv/LogSift/pkg/stats"
	"github.com/sujanv/LogSift/pkg/geodb"
	"github.com/sujanv/LogSift/pkg/query"
	"github.com/sujanv/LogSift/pkg/replay"
	"github.com/sujanv/LogSift/pkg/security"
	"github.com/sujanv/LogSift/pkg/archive"
	"github.com/sujanv/LogSift/pkg/export"
	"github.com/sujanv/LogSift/pkg/session"
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
		case "pipe":
		runPipe(os.Args[2:])
	case "security":
		runSecurity(os.Args[2:])
	case "geo":
		runGeo(os.Args[2:])
	case "replay":
		runReplay(os.Args[2:])
		case "archive", "partition":
		runArchive(os.Args[2:])
	case "session":
		runSession(os.Args[2:])
	case "metrics":
		runMetrics(os.Args[2:])
	case "search":
		runSearch(os.Args[2:])
	case "summarize", "stats":
		runSummarize(os.Args[2:])
	case "index":
		runIndex(os.Args[2:])
	case "query", "sql":
		runQuery(os.Args[2:])
	case "mask":
		runMask(os.Args[2:])
	case "diff":
		runDiff(os.Args[2:])
	case "tail":
		runTail(os.Args[2:])
	case "correlate":
		runCorrelate(os.Args[2:])
	case "serve":
		runServe(os.Args[2:])
	case "html":
		runHTML(os.Args[2:])
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
  tui         Launch interactive Bubble Tea terminal log explorer
  tail        Continuously tail a file like 'tail -f'
  diff        Compare two log files and identify regressions/drift
  mask        Sanitize sensitive PII (emails, cards, tokens, SSNs)
  correlate   Assemble distributed trace waterfall by trace_id
  html        Generate standalone offline HTML dashboard report
  index       Index log file into SQLite for instant SQL querying
  query       Execute arbitrary SQL against indexed SQLite database
  serve       Run embedded HTTP REST API and monitoring server
  help        Show help for LogSift`)
}

func runMask(args []string) {
	fs := flag.NewFlagSet("mask", flag.ExitOnError)
	strat := fs.String("strategy", "redact", "Masking strategy: redact, mask, hash")
	_ = fs.Parse(args)
	files := fs.Args()

	r := getReader(files)
	defer r.Close()

	sr := reader.NewStreamReader(r)
	defer sr.Close()

	m := mask.NewMasker(mask.ParseStrategy(*strat))

	for {
		line, _, _, err := sr.NextLine()
		if err != nil {
			break
		}
		fmt.Println(m.MaskText(string(line)))
	}
}

func runDiff(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: logsift diff <base.log> <target.log>")
		os.Exit(1)
	}
	baseSnap := captureSnapshot(args[0])
	targetSnap := captureSnapshot(args[1])
	report := diff.Compare(baseSnap, targetSnap)
	fmt.Println(report.FormatDelta())
}

func captureSnapshot(filePath string) *diff.Snapshot {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open %s: %v\n", filePath, err)
		os.Exit(1)
	}
	defer f.Close()

	sr := reader.NewStreamReader(f)
	p := parser.NewJSONParser()
	lat := stats.NewLatencyTracker()
	errSum := stats.NewErrorSummarizer()

	var total, errs int64
	rec := model.DefaultRecordPool.Get()
	defer model.DefaultRecordPool.Put(rec)

	for {
		line, _, _, err := sr.NextLine()
		if err != nil {
			break
		}
		total++
		rec.Reset()
		_ = p.Parse(line, rec)
		if rec.Latency > 0 {
			lat.Add(rec.Latency)
		}
		if rec.Level >= model.LevelError || rec.StatusCode >= 500 {
			errs++
			errSum.Add(rec.Message, rec.Timestamp)
		}
	}

	return &diff.Snapshot{
		TotalLogs:  total,
		ErrorCount: errs,
		Latency:    lat.Snapshot(),
		Errors:     errSum.TopClusters(10),
	}
}

func runTail(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: logsift tail <file>")
		os.Exit(1)
	}
	follower := reader.NewTailFollower(args[0], 250*time.Millisecond)
	lines := make(chan []byte, 100)

	go func() {
		_ = follower.Follow(context.Background(), lines)
	}()

	for line := range lines {
		fmt.Println(string(line))
	}
}

func runCorrelate(args []string) {
	fs := flag.NewFlagSet("correlate", flag.ExitOnError)
	traceID := fs.String("trace", "", "Target trace ID")
	_ = fs.Parse(args)

	if *traceID == "" {
		fmt.Fprintln(os.Stderr, "error: --trace flag is required")
		os.Exit(1)
	}

	files := fs.Args()
	r := getReader(files)
	defer r.Close()

	sr := reader.NewStreamReader(r)
	p := parser.NewJSONParser()
	tc := correlate.NewTraceCorrelator()

	rec := model.DefaultRecordPool.Get()
	defer model.DefaultRecordPool.Put(rec)

	for {
		line, _, _, err := sr.NextLine()
		if err != nil {
			break
		}
		rec.Reset()
		_ = p.Parse(line, rec)
		tc.Ingest(rec)
	}

	trace := tc.GetTrace(*traceID)
	fmt.Println(correlate.FormatTraceTimeline(*traceID, trace))
}

func runHTML(args []string) {
	fs := flag.NewFlagSet("html", flag.ExitOnError)
	outFile := fs.String("out", "dashboard.html", "Output HTML path")
	_ = fs.Parse(args)

	files := fs.Args()
	r := getReader(files)
	defer r.Close()

	sr := reader.NewStreamReader(r)
	p := parser.NewJSONParser()
	lat := stats.NewLatencyTracker()
	errSum := stats.NewErrorSummarizer()

	var total int64
	rec := model.DefaultRecordPool.Get()
	defer model.DefaultRecordPool.Put(rec)

	for {
		line, _, _, err := sr.NextLine()
		if err != nil {
			break
		}
		total++
		rec.Reset()
		_ = p.Parse(line, rec)
		if rec.Latency > 0 {
			lat.Add(rec.Latency)
		}
		if rec.Level >= model.LevelError {
			errSum.Add(rec.Message, rec.Timestamp)
		}
	}

	f, err := os.Create(*outFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create HTML file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	gen := output.NewHTMLReportGenerator(f)
	_ = gen.Generate(total, time.Second, lat.Snapshot(), errSum.TopClusters(10), nil)
	fmt.Printf("HTML dashboard generated: %s\n", *outFile)
}

func runServe(args []string) {
	fs := flag.NewFlagSet("serve", flag.ExitOnError)
	port := fs.String("port", "8080", "Server port")
	_ = fs.Parse(args)

	srv := server.NewServer(":"+*port, stats.NewLatencyTracker(), stats.NewErrorSummarizer())
	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
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
	levelFilter := fs.String("level", "", "Filter by log level")
	exprFilter := fs.String("filter", "", "Filter expression")
	sinceFilter := fs.String("since", "", "Filter logs within relative time duration")
	format := fs.String("format", "auto", "Log format")
	outMode := fs.String("output", "terminal", "Output mode")
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
	reportFile := fs.String("report", "", "Write markdown report")
	format := fs.String("format", "auto", "Log format")
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
	var totalRecords, totalErrors int64

	for {
		line, lineNum, offset, err := sr.NextLine()
		if err != nil {
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
			fmt.Fprintf(os.Stderr, "failed to create report: %v\n", err)
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

func runPipe(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: logsift pipe '<pipeline-query>' [file]")
		os.Exit(1)
	}
	stages, err := query.ParsePipeline(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "pipe parse error: %v\n", err)
		os.Exit(1)
	}
	exec, err := query.CompilePipeline(stages)
	if err != nil {
		fmt.Fprintf(os.Stderr, "pipe compile error: %v\n", err)
		os.Exit(1)
	}

	files := args[1:]
	r := getReader(files)
	defer r.Close()

	sr := reader.NewStreamReader(r)
	p := parser.NewJSONParser()
	rec := model.DefaultRecordPool.Get()
	defer model.DefaultRecordPool.Put(rec)

	for {
		line, _, _, err := sr.NextLine()
		if err != nil {
			break
		}
		rec.Reset()
		_ = p.Parse(line, rec)
		if exec.Process(rec) && exec.TopKResult() == nil {
			fmt.Println(string(line))
		}
	}

	if topk := exec.TopKResult(); len(topk) > 0 {
		fmt.Println("\nTop Results:")
		for _, item := range topk {
			fmt.Printf("  %6d %s\n", item.Count, item.Key)
		}
	}
}

func runSecurity(args []string) {
	files := args
	r := getReader(files)
	defer r.Close()

	sr := reader.NewStreamReader(r)
	p := parser.NewJSONParser()
	td := security.NewThreatDetector()
	rec := model.DefaultRecordPool.Get()
	defer model.DefaultRecordPool.Put(rec)

	for {
		line, _, _, err := sr.NextLine()
		if err != nil {
			break
		}
		rec.Reset()
		_ = p.Parse(line, rec)
		td.InspectRecord(rec)
	}

	incidents := td.Incidents()
	fmt.Printf("\n=== LogSift Security Threat Audit (%d incidents) ===\n\n", len(incidents))
	for _, inc := range incidents {
		fmt.Printf("[%s] %s from IP %s: %s (%s)\n",
			inc.Severity, inc.ThreatType, inc.ClientIP, inc.Description, inc.RequestURI)
	}
}

func runGeo(args []string) {
	files := args
	r := getReader(files)
	defer r.Close()

	sr := reader.NewStreamReader(r)
	p := parser.NewJSONParser()
	resolver := geodb.NewResolver()
	counts := make(map[string]int)
	rec := model.DefaultRecordPool.Get()
	defer model.DefaultRecordPool.Put(rec)

	for {
		line, _, _, err := sr.NextLine()
		if err != nil {
			break
		}
		rec.Reset()
		_ = p.Parse(line, rec)
		if rec.IP != "" {
			info := resolver.Lookup(rec.IP)
			counts[info.CountryName]++
		}
	}

	fmt.Println("\n=== Geographic Traffic Breakdown ===")
	for country, cnt := range counts {
		fmt.Printf("  %-30s %d requests\n", country, cnt)
	}
}

func runReplay(args []string) {
	fs := flag.NewFlagSet("replay", flag.ExitOnError)
	target := fs.String("target", "http://localhost:8080", "Target server base URL")
	_ = fs.Parse(args)

	files := fs.Args()
	r := getReader(files)
	defer r.Close()

	sr := reader.NewStreamReader(r)
	p := parser.NewJSONParser()
	replayer := replay.NewTrafficReplayer(*target, 1.0, 4)
	rec := model.DefaultRecordPool.Get()
	defer model.DefaultRecordPool.Put(rec)

	var count int
	for {
		line, _, _, err := sr.NextLine()
		if err != nil {
			break
		}
		rec.Reset()
		_ = p.Parse(line, rec)
		if rec.Path != "" {
			code, dur, _ := replayer.ReplaySingle(context.Background(), rec)
			fmt.Printf("Replayed %s %s -> HTTP %d (%v)\n", rec.Method, rec.Path, code, dur)
			count++
		}
	}
	fmt.Printf("Replayed %d requests against %s\n", count, *target)
}

func runArchive(args []string) {
	fs := flag.NewFlagSet("archive", flag.ExitOnError)
	outDir := fs.String("out", "archive", "Destination directory")
	mode := fs.String("by", "date", "Partition mode: date, hour, level")
	_ = fs.Parse(args)

	files := fs.Args()
	r := getReader(files)
	defer r.Close()

	sr := reader.NewStreamReader(r)
	p := parser.NewJSONParser()
	arch := archive.NewChunkArchiver(*outDir, archive.PartitionMode(*mode))
	defer arch.Close()

	rec := model.DefaultRecordPool.Get()
	defer model.DefaultRecordPool.Put(rec)

	for {
		line, _, _, err := sr.NextLine()
		if err != nil {
			break
		}
		rec.Reset()
		_ = p.Parse(line, rec)
		rec.Raw = line
		_ = arch.Archive(rec)
	}
	fmt.Printf("Archived logs into %s (partitioned by %s)\n", *outDir, *mode)
}

func runSession(args []string) {
	files := args
	r := getReader(files)
	defer r.Close()

	sr := reader.NewStreamReader(r)
	p := parser.NewJSONParser()
	st := session.NewSessionTracker(30 * time.Minute)
	rec := model.DefaultRecordPool.Get()
	defer model.DefaultRecordPool.Put(rec)

	for {
		line, _, _, err := sr.NextLine()
		if err != nil {
			break
		}
		rec.Reset()
		_ = p.Parse(line, rec)
		st.Ingest(rec)
	}

	sessions := st.CompleteFinalizes()
	analytics := session.AnalyzeSessions(sessions)

	fmt.Printf("\n=== User Session Analysis (%d sessions) ===\n", analytics.TotalSessions)
	fmt.Printf("Bounce Rate:     %.2f%%\n", analytics.BounceRate)
	fmt.Printf("Avg Requests:    %.1f per session\n", analytics.AvgRequests)
	fmt.Printf("Avg Duration:    %v\n", analytics.AvgDuration)
}

func runMetrics(args []string) {
	files := args
	r := getReader(files)
	defer r.Close()

	sr := reader.NewStreamReader(r)
	p := parser.NewJSONParser()
	lat := stats.NewLatencyTracker()
	var total, errs int64
	rec := model.DefaultRecordPool.Get()
	defer model.DefaultRecordPool.Put(rec)

	for {
		line, _, _, err := sr.NextLine()
		if err != nil {
			break
		}
		total++
		rec.Reset()
		_ = p.Parse(line, rec)
		if rec.Latency > 0 {
			lat.Add(rec.Latency)
		}
		if rec.Level >= model.LevelError {
			errs++
		}
	}

	prom := export.FormatPrometheusMetrics(total, errs, lat.Snapshot(), nil)
	fmt.Print(prom)
}
