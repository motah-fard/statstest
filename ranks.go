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

// rankValues assigns average ranks to the values in x, preserving x's order.
func rankValues(x []float64) []float64 {
	type indexedValue struct {
		value float64
		index int
	}

	items := make([]indexedValue, len(x))
	for i, v := range x {
		items[i] = indexedValue{value: v, index: i}
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].value < items[j].value
	})

	ranks := make([]float64, len(x))

	i := 0
	for i < len(items) {
		j := i + 1
		for j < len(items) && items[j].value == items[i].value {
			j++
		}

		avgRank := (float64(i+1) + float64(j)) / 2.0
		for k := i; k < j; k++ {
			ranks[items[k].index] = avgRank
		}

		i = j
	}

	return ranks
}

// rankGroups pools values across all groups, assigns average ranks to ties,
// and returns the per-group ranks (preserving each group's order), the tie
// counts across the pooled data, and the total pooled sample size.
func rankGroups(groups [][]float64) (ranks [][]float64, tieCounts []int, n int) {
	type item struct {
		value      float64
		group, idx int
	}

	var all []item
	for gi, g := range groups {
		for idx, v := range g {
			all = append(all, item{value: v, group: gi, idx: idx})
		}
	}
	n = len(all)

	sort.Slice(all, func(i, j int) bool {
		return all[i].value < all[j].value
	})

	ranks = make([][]float64, len(groups))
	for gi, g := range groups {
		ranks[gi] = make([]float64, len(g))
	}

	i := 0
	for i < len(all) {
		j := i + 1
		for j < len(all) && all[j].value == all[i].value {
			j++
		}

		avgRank := (float64(i+1) + float64(j)) / 2.0
		for k := i; k < j; k++ {
			ranks[all[k].group][all[k].idx] = avgRank
		}

		if j-i > 1 {
			tieCounts = append(tieCounts, j-i)
		}

		i = j
	}

	return ranks, tieCounts, n
}
