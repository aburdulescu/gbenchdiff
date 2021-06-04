package main

import (
	"fmt"
	"math"
	"os"
	"text/tabwriter"
)

type Sample struct {
	Name    string
	Values  []float64
	RValues []float64 // without outliers
	Min     float64
	Mean    float64
	Max     float64
}

func (s Sample) String() string {
	return fmt.Sprintf("%s: %.4f %.4f %.4f", s.Name, s.Min, s.Mean, s.Max)
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

func (o Sample) Diff(n Sample) {
	w := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)

	fmt.Fprintln(w, "Benchmark\tDelta\tOld\tNew")
	fmt.Fprintln(w, "---------\t-----\t---\t---")

	diff := ((o.Mean - n.Mean) / math.Abs(o.Mean)) * 100

	fmt.Fprintf(w, "%s", o.Name)

	if diff > 0 {
		fmt.Fprintf(w, "\t+%.2f%%", diff)
	} else {
		fmt.Fprintf(w, "\t%.2f%%", diff)
	}
	fmt.Fprintf(w, "\t%.2f\t%.2f\n", o.Mean, n.Mean)

	w.Flush()
}
