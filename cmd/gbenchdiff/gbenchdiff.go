package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"text/tabwriter"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("need two args: old.json and new.json")
	}

	oldFile, err := os.Open(os.Args[1])
	if err != nil {
		return err
	}
	defer oldFile.Close()
	newFile, err := os.Open(os.Args[2])
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

	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)

	fmt.Fprintln(w, "Benchmark\tTime(%)\tCPU(%)\tTime old\tTime new\tCPU old\tCPU new")
	fmt.Fprintln(w, "---------\t----\t---\t--------\t--------\t-------\t-------")

	for _, benchmark := range oldRes.Benchmarks {
		i := findBenchmarkByName(newRes.Benchmarks, benchmark.Name)
		if i == -1 {
			continue
		}
		benchmark.PrintDiff(w, newRes.Benchmarks[i])
	}

	w.Flush()

	return nil
}

func (old Benchmark) PrintDiff(w io.Writer, new Benchmark) {
	realTimeDiff := ((new.RealTime - old.RealTime) / math.Abs(old.RealTime)) * 100
	cpuTimeDiff := ((new.CPUTime - old.CPUTime) / math.Abs(old.CPUTime)) * 100
	fmt.Fprintf(w, "%s", old.Name)
	if realTimeDiff > 0 {
		fmt.Fprintf(w, "\t+%.2f", realTimeDiff)
	} else {
		fmt.Fprintf(w, "\t%.2f", realTimeDiff)
	}
	if cpuTimeDiff > 0 {
		fmt.Fprintf(w, "\t+%.2f", cpuTimeDiff)
	} else {
		fmt.Fprintf(w, "\t%.2f", cpuTimeDiff)
	}
	fmt.Fprintf(w, "\t%.2f\t%.2f\t%.2f\t%.2f\n", old.RealTime, new.RealTime, old.CPUTime, new.CPUTime)

}

func findBenchmarkByName(benchmarks []Benchmark, name string) int {
	for i, b := range benchmarks {
		if b.Name == name {
			return i
		}
	}
	return -1
}

type Result struct {
	Benchmarks []Benchmark `json:"benchmarks"`
}

// TODO: take in consideration aggregates if present
type Benchmark struct {
	Name            string  `json:"name"`
	RunName         string  `json:"run_name"`
	RunType         string  `json:"run_type"`
	Repetitions     uint64  `json:"repetitions"`
	RepetitionIndex uint64  `json:"repetition_index"`
	Threads         int     `json:"threads"`
	Iterations      uint64  `json:"iterations"`
	RealTime        float64 `json:"real_time"`
	CPUTime         float64 `json:"cpu_time"`
	TimeUnit        string  `json:"time_unit"`
}
