package statstest

import (
	"fmt"
	"sort"
)

type pvalIndex struct {
	value float64
	index int
}

// AdjustPValues adjusts a slice of p-values using the requested
// multiple-testing correction method.
//
// Supported methods are Bonferroni, Holm, and Benjamini-Hochberg.
// The returned slice preserves the original order of the input p-values.

func AdjustPValues(p []float64, method PAdjustMethod) ([]float64, error) {
	if len(p) == 0 {
		return nil, ErrEmptySample
	}

	for _, v := range p {
		if v < 0 || v > 1 {
			return nil, ErrInvalidPValue
		}
	}

	switch method {
	case Bonferroni:
		return adjustBonferroni(p), nil
	case Holm:
		return adjustHolm(p), nil
	case BenjaminiHochberg:
		return adjustBH(p), nil
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidPAdjustMethod, method)
	}
}

func adjustBonferroni(p []float64) []float64 {
	n := float64(len(p))
	out := make([]float64, len(p))
	for i, v := range p {
		adj := v * n
		if adj > 1 {
			adj = 1
		}
		out[i] = adj
	}
	return out
}

func adjustHolm(p []float64) []float64 {
	n := len(p)
	indexed := make([]pvalIndex, n)
	for i, v := range p {
		indexed[i] = pvalIndex{value: v, index: i}
	}

	sort.Slice(indexed, func(i, j int) bool {
		return indexed[i].value < indexed[j].value
	})

	adjustedOrdered := make([]float64, n)
	for i := 0; i < n; i++ {
		adjustedOrdered[i] = float64(n-i) * indexed[i].value
	}

	for i := 1; i < n; i++ {
		if adjustedOrdered[i] < adjustedOrdered[i-1] {
			adjustedOrdered[i] = adjustedOrdered[i-1]
		}
	}

	out := make([]float64, n)
	for i := 0; i < n; i++ {
		v := adjustedOrdered[i]
		if v > 1 {
			v = 1
		}
		out[indexed[i].index] = v
	}

	return out
}

func adjustBH(p []float64) []float64 {
	n := len(p)
	indexed := make([]pvalIndex, n)
	for i, v := range p {
		indexed[i] = pvalIndex{value: v, index: i}
	}

	sort.Slice(indexed, func(i, j int) bool {
		return indexed[i].value < indexed[j].value
	})

	adjustedOrdered := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		rank := float64(i + 1)
		v := indexed[i].value * float64(n) / rank
		if i == n-1 {
			adjustedOrdered[i] = v
		} else {
			if v > adjustedOrdered[i+1] {
				v = adjustedOrdered[i+1]
			}
			adjustedOrdered[i] = v
		}
	}

	out := make([]float64, n)
	for i := 0; i < n; i++ {
		v := adjustedOrdered[i]
		if v > 1 {
			v = 1
		}
		out[indexed[i].index] = v
	}

	return out
}
