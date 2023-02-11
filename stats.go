package main

import "math"

// Percentile returns the pctileth value from the sample. This uses
// interpolation method R8 from Hyndman and Fan (1996).
//
// pctile will be capped to the range [0, 1]. If len(xs) == 0,
// returns NaN.
//
// Percentile(0.5) is the median. Percentile(0.25) and
// Percentile(0.75) are the first and third quartiles, respectively.
func Percentile(xs []float64, pctile float64) float64 {
	if len(xs) == 0 {
		return math.NaN()
	}
	N := float64(len(xs))
	// n := pctile * (N + 1) // R6
	n := 1/3.0 + pctile*(N+1/3.0) // R8
	kf, frac := math.Modf(n)
	k := int(kf)
	if k <= 0 {
		return xs[0]
	} else if k >= len(xs) {
		return xs[len(xs)-1]
	}
	return xs[k-1] + frac*(xs[k]-xs[k-1])
}

// Mean returns the arithmetic mean of xs.
func Mean(xs []float64) float64 {
	if len(xs) == 0 {
		return math.NaN()
	}
	m := 0.0
	for i, x := range xs {
		m += (x - m) / float64(i+1)
	}
	return m
}

// Bounds returns the minimum and maximum values of xs.
func Bounds(xs []float64) (min float64, max float64) {
	if len(xs) == 0 {
		return math.NaN(), math.NaN()
	}
	min, max = xs[0], xs[0]
	for _, x := range xs {
		if x < min {
			min = x
		}
		if x > max {
			max = x
		}
	}
	return
}
