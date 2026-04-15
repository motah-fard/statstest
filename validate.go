package statstest

import (
	"math"
)

// FloatSlice checks whether a float slice is non-empty, meets the minimum
// length, and contains no NaN or Inf values.
func floatSlice(x []float64, minLen int) error {
	if len(x) == 0 {
		return ErrEmptySample
	}
	if len(x) < minLen {
		return ErrSampleTooSmall
	}
	for _, v := range x {
		if math.IsNaN(v) {
			return ErrContainsNaN
		}
		if math.IsInf(v, 0) {
			return ErrContainsInf
		}
	}
	return nil
}

// Confidence checks whether conf is in (0, 1).
func confidence(conf float64) error {
	if conf <= 0 || conf >= 1 {
		return ErrInvalidConfidenceLevel
	}
	return nil
}

// PairedLengths checks whether two slices have the same length.
func pairedLengths[T any](x, y []T) error {
	if len(x) != len(y) {
		return ErrMismatchedLengths
	}
	return nil
}

// Alternative checks whether alt is one of the supported string values.
func alternative1(alt string) error {
	switch alt {
	case "two-sided", "less", "greater":
		return nil
	default:
		return ErrInvalidAlternative
	}
}

// ContingencyTable checks whether a contingency table is rectangular,
// at least 2x2, and contains no negative counts.
func contingencyTable(observed [][]int) error {
	if len(observed) < 2 {
		return ErrInvalidTable
	}
	if len(observed[0]) < 2 {
		return ErrInvalidTable
	}

	cols := len(observed[0])
	for i := range observed {
		if len(observed[i]) != cols {
			return ErrInvalidTable
		}
		for _, v := range observed[i] {
			if v < 0 {
				return ErrNegativeCount
			}
		}
	}
	return nil
}

// Table2x2 checks whether all counts in a 2x2 table are non-negative and
// whether the grand total is positive.
func table2x2(table [2][2]int) error {
	total := 0
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			if table[i][j] < 0 {
				return ErrNegativeCount
			}
			total += table[i][j]
		}
	}
	if total == 0 {
		return ErrInvalidTable
	}
	return nil
}
