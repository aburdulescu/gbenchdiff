package main

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"sort"

	"bandr.me/p/gbenchdiff/internal/stats"
)

const alpha = 0.05

type Metric struct {
	Name     string
	TimeUnit string
	RealTime Sample
	CPUTime  Sample
}

type Sample struct {
	Values  []float64
	RValues []float64 // without outliers
	Min     float64
	Mean    float64
	Max     float64
}

func (s *Sample) removeOutliers() {
	q1 := Percentile(s.Values, 0.25)
	q3 := Percentile(s.Values, 0.75)
	lo := q1 - 1.5*(q3-q1)
	hi := q3 + 1.5*(q3-q1)
	for _, value := range s.Values {
		if value >= lo && value <= hi {
			s.RValues = append(s.RValues, value)
		}
	}
}

func (s *Sample) ComputeStats() {
	s.removeOutliers()
	s.Min, s.Max = Bounds(s.RValues)
	s.Mean = Mean(s.RValues)
}

func (o Sample) Print(w io.Writer, n Sample, tu string) {
	u, err := stats.MannWhitneyUTest(o.RValues, n.RValues, stats.LocationDiffers)

	pval := u.P

	delta := "~"
	note := ""

	switch {
	case errors.Is(err, stats.ErrZeroVariance):
		note = "(zero variance)"
	case errors.Is(err, stats.ErrSampleSize):
		note = "(too few samples)"
	case errors.Is(err, stats.ErrSamplesEqual):
		note = "(all equal)"
	case err != nil:
		note = fmt.Sprintf("(%s)", err)
	case pval < alpha:
		if n.Mean == o.Mean {
			delta = "0.00%"
		} else {
			pct := ((n.Mean - o.Mean) / o.Mean) * 100.0
			delta = fmt.Sprintf("%+.2f%%", pct)
		}
	}

	if note == "" && pval != -1 {
		note = fmt.Sprintf("(p=%0.2f n=%d+%d)", pval, len(o.RValues), len(n.RValues))
	}

	fmt.Fprintf(w, "\t%s\t%s", delta, note)
	fmt.Fprintf(w, "\t%.2f%s\t%.2f%s", o.Mean, tu, n.Mean, tu)
}

func findMetric(m []Metric, name string) int {
	for i := range m {
		if m[i].Name == name {
			return i
		}
	}
	return -1
}

func GetMetrics(benchmarks []Benchmark, filterRe *regexp.Regexp) []Metric {
	var metrics []Metric
	for _, b := range benchmarks {
		if filterRe != nil && !filterRe.MatchString(b.Name) {
			continue
		}
		if b.RunType != "iteration" {
			continue
		}
		i := findMetric(metrics, b.Name)
		if i == -1 {
			metrics = append(metrics, Metric{
				Name:     b.Name,
				TimeUnit: b.TimeUnit,
			})
			i = len(metrics) - 1
		}
		metrics[i].RealTime.Values = append(metrics[i].RealTime.Values, b.RealTime)
		metrics[i].CPUTime.Values = append(metrics[i].CPUTime.Values, b.CPUTime)
	}
	for i := range metrics {
		r := metrics[i].RealTime.Values
		sort.Float64s(r)
		metrics[i].RealTime.Values = r
		metrics[i].RealTime.ComputeStats()

		c := metrics[i].CPUTime.Values
		sort.Float64s(c)
		metrics[i].CPUTime.Values = c
		metrics[i].CPUTime.ComputeStats()
	}
	return metrics
}
