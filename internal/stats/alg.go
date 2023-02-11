// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stats

// Miscellaneous helper algorithms

func maxint(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minint(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sumint(xs []int) int {
	sum := 0
	for _, x := range xs {
		sum += x
	}
	return sum
}
