package main

import (
	"fmt"
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
