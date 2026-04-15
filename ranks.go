package statstest

import "sort"

type rankedValue struct {
	value float64
	fromX bool
	rank  float64
}

// rankWithAverageTies combines x and y, sorts the pooled values,
// and assigns average ranks to tied observations.
func rankWithAverageTies(x, y []float64) ([]rankedValue, []int) {
	all := make([]rankedValue, 0, len(x)+len(y))

	for _, v := range x {
		all = append(all, rankedValue{value: v, fromX: true})
	}
	for _, v := range y {
		all = append(all, rankedValue{value: v, fromX: false})
	}

	sort.Slice(all, func(i, j int) bool {
		return all[i].value < all[j].value
	})

	var tieCounts []int

	i := 0
	for i < len(all) {
		j := i + 1
		for j < len(all) && all[j].value == all[i].value {
			j++
		}

		startRank := float64(i + 1)
		endRank := float64(j)
		avgRank := (startRank + endRank) / 2.0

		for k := i; k < j; k++ {
			all[k].rank = avgRank
		}

		if j-i > 1 {
			tieCounts = append(tieCounts, j-i)
		}

		i = j
	}

	return all, tieCounts
}

// sumRanksForX returns the sum of pooled ranks assigned to sample x.
func sumRanksForX(ranked []rankedValue) float64 {
	var sum float64
	for _, rv := range ranked {
		if rv.fromX {
			sum += rv.rank
		}
	}
	return sum
}
