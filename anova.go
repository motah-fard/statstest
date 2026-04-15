package statstest

import (
	"gonum.org/v1/gonum/stat/distuv"
)

// OneWayANOVA performs a one-way analysis of variance across two or more groups.
//
// The null hypothesis is that all group means are equal.
func OneWayANOVA(groups ...[]float64) (OneWayANOVAResult, error) {
	if len(groups) < 2 {
		return OneWayANOVAResult{}, ErrTooFewGroups
	}

	totalN := 0
	var grandSum float64

	for _, g := range groups {
		if err := validateSample(g, 2); err != nil {
			return OneWayANOVAResult{}, err
		}
		totalN += len(g)
		for _, v := range g {
			grandSum += v
		}
	}

	k := len(groups)
	if totalN <= k {
		return OneWayANOVAResult{}, ErrInsufficientDF
	}

	grandMean := grandSum / float64(totalN)

	var ssBetween float64
	var ssWithin float64
	var ssTotal float64

	for _, g := range groups {
		groupMean := mean(g)
		n := float64(len(g))

		diffGroup := groupMean - grandMean
		ssBetween += n * diffGroup * diffGroup

		for _, v := range g {
			diffWithin := v - groupMean
			ssWithin += diffWithin * diffWithin

			diffTotal := v - grandMean
			ssTotal += diffTotal * diffTotal
		}
	}

	dfBetween := k - 1
	dfWithin := totalN - k
	if dfBetween <= 0 || dfWithin <= 0 {
		return OneWayANOVAResult{}, ErrInsufficientDF
	}

	msBetween := ssBetween / float64(dfBetween)
	msWithin := ssWithin / float64(dfWithin)
	if msWithin <= 0 {
		return OneWayANOVAResult{}, ErrZeroVariance
	}

	f := msBetween / msWithin

	dist := distuv.F{
		D1: float64(dfBetween),
		D2: float64(dfWithin),
	}
	p := 1 - dist.CDF(f)

	etaSquared := 0.0
	if ssTotal > 0 {
		etaSquared = ssBetween / ssTotal
	}

	return OneWayANOVAResult{
		FStatistic: f,
		PValue:     p,
		DFBetween:  dfBetween,
		DFWithin:   dfWithin,
		SSTotal:    ssTotal,
		SSBetween:  ssBetween,
		SSWithin:   ssWithin,
		EtaSquared: etaSquared,
		Method:     "One-way ANOVA",
	}, nil
}
