package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
)

const usageExtra = `
For each benchmark in both files, the tool will:
- remove outliers with interquartile range rule
- perform significance test(Man-Whitney U-test)
- print % change in mean from the first to the second file
- print the p-value and sample sizes from a test of the two distributions of benchmark times

Small p-values indicate that the two distributions are significantly different.
If the test indicates that there was no significant change between the two
benchmarks (defined as p > 0.05), a single ~ will be displayed instead of
the percent change.

Notes:
- run the benchmark with --benchmark_repetitions(=10 should be enough)
- run the benchmark with --benchmark_out=file.json
`

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options] old.json new.json\n", os.Args[0])
	fmt.Fprint(os.Stderr, "options:\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, usageExtra)
	os.Exit(1)
}

func run() error {
	// var fHtml bool
	var fNoCtxCheck bool
	// flag.BoolVar(&fHtml, "html", false, "print result as HTML")
	flag.BoolVar(&fNoCtxCheck, "no-ctx", false, "don't compare benchmark contexts")
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		usage()
	}

	oldFilepath := args[0]
	newFilepath := args[1]

	oldFile, err := os.Open(oldFilepath)
	if err != nil {
		return err
	}
	defer oldFile.Close()
	newFile, err := os.Open(newFilepath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	var oldRes Result
	if err := json.NewDecoder(oldFile).Decode(&oldRes); err != nil {
		return err
	}
	var newRes Result
	if err := json.NewDecoder(newFile).Decode(&newRes); err != nil {
		return err
	}

	if !fNoCtxCheck {
		if err := oldRes.Context.Equals(newRes.Context); err != nil {
			return fmt.Errorf("context check failed: %v", err)
		}
	}

	metrics_o := GetMetrics(oldRes.Benchmarks)
	metrics_n := GetMetrics(newRes.Benchmarks)

	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)

	fmt.Fprintln(w, "name\treal\tnote\told\tnew\tcpu\tnote\told\tnew")
	fmt.Fprintln(w, "----\t----\t----\t---\t---\t---\t----\t---\t---")

	for _, m_o := range metrics_o {
		i := findMetric(metrics_n, m_o.Name)
		if i == -1 {
			continue
		}
		m_n := metrics_n[i]

		fmt.Fprintf(w, "%s", m_o.Name)

		m_o.RealTime.Print(w, m_n.RealTime)
		m_o.CPUTime.Print(w, m_n.CPUTime)

		fmt.Fprintln(w)
	}

	w.Flush()

	return nil
}
