package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
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
	var fWithCPUTime bool

	// flag.BoolVar(&fHtml, "html", false, "print result as HTML")
	flag.BoolVar(&fNoCtxCheck, "no-ctx", false, "don't compare benchmark contexts")
	flag.BoolVar(&fWithCPUTime, "with-cpu", false, "compare also CPU time")

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

	oldMetrics := GetMetrics(oldRes.Benchmarks)
	newMetrics := GetMetrics(newRes.Benchmarks)

	printer := Printer{
		w: tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0),
	}

	if err := printer.Print("real", oldMetrics, newMetrics); err != nil {
		return err
	}

	if !fWithCPUTime {
		return nil
	}

	return printer.Print("cpu", oldMetrics, newMetrics)
}

type Printer struct {
	w *tabwriter.Writer
}

func (p Printer) Print(what string, old, new []Metric) error {
	if what != "real" && what != "cpu" {
		return fmt.Errorf("unknown what value '%s'", what)
	}

	fmt.Fprintf(p.w, "name\t%s\tnote\told\tnew\n", what)
	fmt.Fprintf(p.w, "----\t%s\t----\t---\t---\n", strings.Repeat("-", len(what)))

	for _, o := range old {
		i := findMetric(new, o.Name)
		if i == -1 {
			continue
		}
		n := new[i]

		if n.TimeUnit != o.TimeUnit {
			return fmt.Errorf(
				"benchmarks have different time units: old=%s, new=%s",
				o.TimeUnit, n.TimeUnit)
		}

		fmt.Fprintf(p.w, "%s", n.Name)

		if what == "real" {
			o.RealTime.Print(p.w, n.RealTime, n.TimeUnit)
		} else {
			o.CPUTime.Print(p.w, n.CPUTime, n.TimeUnit)
		}

		fmt.Fprintln(p.w)
	}

	return p.w.Flush()
}
