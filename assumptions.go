package statstest

// HasMinimumSampleSize reports whether x has at least min observations.
func HasMinimumSampleSize(x []float64, min int) bool {
	return len(x) >= min
}

// HasPairedEqualLength reports whether two paired samples have equal length.
func HasPairedEqualLength(x, y []float64) bool {
	return len(x) == len(y)
}

// HasAtLeastTwoGroups reports whether at least two groups are provided.
func HasAtLeastTwoGroups(groups ...[]float64) bool {
	return len(groups) >= 2
}

// HasRectangularShape reports whether a contingency table is rectangular.
func HasRectangularShape(table [][]int) bool {
	if len(table) == 0 {
		return false
	}
	cols := len(table[0])
	if cols == 0 {
		return false
	}
	for i := range table {
		if len(table[i]) != cols {
			return false
		}
	}
	return true
}
