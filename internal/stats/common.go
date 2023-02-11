package stats

import (
	"errors"
	"math"
)

var inf = math.Inf(1)
var nan = math.NaN()

var (
	ErrSamplesEqual      = errors.New("all samples are equal")
	ErrSampleSize        = errors.New("sample is too small")
	ErrZeroVariance      = errors.New("sample has zero variance")
	ErrMismatchedSamples = errors.New("samples have different lengths")
)
