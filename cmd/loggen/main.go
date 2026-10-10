package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/sujanv/LogSift/pkg/generator"
)

func main() {
	count := flag.Int64("n", 100000, "Number of log records to generate")
	outFile := flag.String("o", "", "Output file path (default stdout)")
	flag.Parse()

	out := os.Stdout
	if *outFile != "" {
		f, err := os.Create(*outFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create output file: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		out = f
	}

	gen := generator.NewGenerator(42)
	t0 := time.Now()
	if err := gen.GenerateJSONLines(out, *count); err != nil {
		fmt.Fprintf(os.Stderr, "generation error: %v\n", err)
		os.Exit(1)
	}

	elapsed := time.Since(t0)
	rate := float64(*count) / elapsed.Seconds()
	if *outFile != "" {
		fmt.Fprintf(os.Stderr, "Generated %d logs in %v (%.0f logs/sec) -> %s\n", *count, elapsed, rate, *outFile)
	}
}
