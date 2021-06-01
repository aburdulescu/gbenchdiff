package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
	"text/tabwriter"
)

// must run benchmarks with multiple repetitions(at least 10)

// useful links
// https://pkg.go.dev/golang.org/x/perf/cmd/benchstat
// https://pkg.go.dev/golang.org/x/perf@v0.0.0-20210220033136-40a54f11e909/internal/stats
// https://pkg.go.dev/golang.org/x/perf@v0.0.0-20210220033136-40a54f11e909/benchstat

// TODO: error if too few repetitions/samples
// TODO: implement ManWhitneyUTest(if it fails ignore mean difference), needs multiple samples(at least 10)
// TODO: use interquartile range rule to remove outliers

const usage = `gbenchdiff [options] old.json new.json
- remove outliers with interquartile range rule
- perform significance test(Man-Whitney U-test)
- print p-value
- print % difference of mean value
- print times for each file
`

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

type Result struct {
	Benchmarks []Benchmark `json:"benchmarks"`
}

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

	oldMeans := getMeans(oldRes.Benchmarks)
	newMeans := getMeans(newRes.Benchmarks)

	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)

	fmt.Fprintln(w, "Benchmark\tTime(%)\tCPU(%)\tTime old\tTime new\tCPU old\tCPU new")
	fmt.Fprintln(w, "---------\t----\t---\t--------\t--------\t-------\t-------")

	for name, oldMetric := range oldMeans {
		newMetric, ok := newMeans[name]
		if !ok {
			continue
		}
		if newMetric.TimeUnit != oldMetric.TimeUnit {
			return fmt.Errorf("benchmarks %s has different time unit: old=%s, new=%s\n",
				name, oldMetric.TimeUnit, newMetric.TimeUnit)
		}
		oldMetric.PrintDiff(name, newMetric, w)
	}

	w.Flush()

	return nil
}

type Metric struct {
	RealTime float64
	CPUTime  float64
	TimeUnit string
}

func (old Metric) PrintDiff(name string, new Metric, w io.Writer) {
	realTimeDiff := ((old.RealTime - new.RealTime) / math.Abs(old.RealTime)) * 100
	cpuTimeDiff := ((old.CPUTime - new.CPUTime) / math.Abs(old.CPUTime)) * 100

	fmt.Fprintf(w, "%s", name)

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

func getMeans(benchmarks []Benchmark) map[string]Metric {
	means := make(map[string]Metric)
	for _, benchmark := range benchmarks {
		if strings.HasSuffix(benchmark.Name, "_mean") {
			name := strings.TrimSuffix(benchmark.Name, "_mean")
			means[name] = Metric{
				RealTime: benchmark.RealTime,
				CPUTime:  benchmark.CPUTime,
				TimeUnit: benchmark.TimeUnit,
			}
		}
	}
	return means
}
