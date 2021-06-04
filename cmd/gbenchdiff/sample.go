package main

import (
	"fmt"
	"io"
	"math"
	"sort"
	"strings"
)

type Metrics struct {
	Name     string
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

func findMetric(m []Metrics, name string) int {
	for i := range m {
		if m[i].Name == name {
			return i
		}
	}
	return -1
}

func GetMetrics(benchmarks []Benchmark) []Metrics {
	var metrics []Metrics
	for _, b := range benchmarks {
		if strings.HasSuffix(b.Name, "_mean") ||
			strings.HasSuffix(b.Name, "_median") ||
			strings.HasSuffix(b.Name, "_stddev") {
			continue
		}
		i := findMetric(metrics, b.Name)
		if i == -1 {
			metrics = append(metrics, Metrics{Name: b.Name})
			i = len(metrics) - 1
		}
		metrics[i].RealTime.Values = append(metrics[i].RealTime.Values, b.RealTime)
		metrics[i].CPUTime.Values = append(metrics[i].RealTime.Values, b.CPUTime)
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

func (o Sample) Print(w io.Writer, n Sample) {
	diff := ((o.Mean - n.Mean) / math.Abs(o.Mean)) * 100
	if diff > 0 {
		fmt.Fprintf(w, "\t+%.2f%%", diff)
	} else {
		fmt.Fprintf(w, "\t%.2f%%", diff)
	}
	fmt.Fprintf(w, "\t%.2f\t%.2f", o.Mean, n.Mean)
}
