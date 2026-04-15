package statstest

import "sort"

type signedRank struct {
	absDiff float64
	sign    int
	rank    float64
}

func signedRanksFromDifferences(diffs []float64) ([]signedRank, []int) {
	items := make([]signedRank, 0, len(diffs))
	for _, d := range diffs {
		if d == 0 {
			continue
		}
		sign := 1
		if d < 0 {
			sign = -1
			d = -d
		}
		items = append(items, signedRank{
			absDiff: d,
			sign:    sign,
		})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].absDiff < items[j].absDiff
	})

	var tieCounts []int

	i := 0
	for i < len(items) {
		j := i + 1
		for j < len(items) && items[j].absDiff == items[i].absDiff {
			j++
		}

		startRank := float64(i + 1)
		endRank := float64(j)
		avgRank := (startRank + endRank) / 2.0

		for k := i; k < j; k++ {
			items[k].rank = avgRank
		}

		if j-i > 1 {
			tieCounts = append(tieCounts, j-i)
		}

		i = j
	}

	return items, tieCounts
}

func sumPositiveNegativeRanks(items []signedRank) (wPlus, wMinus float64) {
	for _, it := range items {
		if it.sign > 0 {
			wPlus += it.rank
		} else {
			wMinus += it.rank
		}
	}
	return wPlus, wMinus
}
