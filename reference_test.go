// -----------------------------------------------------------------------------
// Reference-value tests
//
// These tests verify results against external statistical references or
// package-documented conventions, rather than only checking that functions run.
//
// Convention notes:
// - Mann-Whitney expected values follow this package's current asymptotic
//   p-value convention.
// - WilcoxonSignedRank reports W as the sum of positive ranks (W+).
// - OneWayANOVA reports FStatistic, DFBetween, and DFWithin.
// - PearsonCorrelation, SpearmanCorrelation, KruskalWallis, and the
//   proportion tests are checked against scipy.stats and
//   statsmodels.stats.proportion (SciPy 1.18.1, statsmodels 0.15.0).
// - ChiSquareGoodnessOfFit and LevenesTest are checked against
//   scipy.stats.chisquare and scipy.stats.levene(center='median').
// - TukeyHSD is checked against scipy.stats.tukey_hsd, and its underlying
//   studentized range distribution against scipy.stats.studentized_range.
// - DunnTest is checked against scikit_posthocs.posthoc_dunn.
// - ShapiroWilk is checked against scipy.stats.shapiro.
// - The power/sample-size functions are checked against
//   statsmodels.stats.power.NormalIndPower, which uses the same
//   normal-approximation formula implemented here.
// -----------------------------------------------------------------------------

package statstest

import (
	"math"
	"testing"
)

func TestReferencePAdjustBH(t *testing.T) {
	p := []float64{0.001, 0.01, 0.03, 0.2}

	got, err := AdjustPValues(p, BenjaminiHochberg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []float64{0.004, 0.02, 0.04, 0.2}
	requireSliceAlmostEqual(t, "BH adjusted p-values", got, want, 1e-12)
}

func almostEqual(got, want, tol float64) bool {
	return math.Abs(got-want) <= tol
}

func requireAlmostEqual(t *testing.T, name string, got, want, tol float64) {
	t.Helper()
	if !almostEqual(got, want, tol) {
		t.Fatalf("%s: got %.12g, want %.12g (tol=%.12g)", name, got, want, tol)
	}
}

func requireSliceAlmostEqual(t *testing.T, name string, got, want []float64, tol float64) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("%s: length mismatch: got %d, want %d", name, len(got), len(want))
	}
	for i := range got {
		if !almostEqual(got[i], want[i], tol) {
			t.Fatalf("%s[%d]: got %.12g, want %.12g (tol=%.12g)", name, i, got[i], want[i], tol)
		}
	}
}

func TestReferenceTTestOneSample(t *testing.T) {
	x := []float64{2.3, 2.5, 2.1, 2.7, 2.4}

	wantStat := 4.000000000000003
	wantP := 0.016130089900092518
	wantDF := 4.0
	tol := 1e-12

	res, err := TTestOneSample(x, 2.0, TwoSided, 0.95)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "statistic", res.Statistic, wantStat, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
	requireAlmostEqual(t, "df", res.DF, wantDF, tol)
}

func TestReferenceTTestWelch(t *testing.T) {
	x := []float64{12.1, 11.8, 12.5, 12.0, 11.9}
	y := []float64{10.2, 10.4, 10.1, 10.3}

	wantStat := 13.21250013653373
	wantP := 0.000012203818575066068
	wantDF := 5.96150036361814
	tol := 1e-10

	res, err := TTestTwoSample(x, y, TwoSampleTOptions{
		Alternative:     TwoSided,
		ConfidenceLevel: 0.95,
		EqualVariance:   false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "statistic", res.Statistic, wantStat, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
	requireAlmostEqual(t, "df", res.DF, wantDF, tol)
}

func TestReferenceTTestPaired(t *testing.T) {
	x := []float64{10.2, 9.8, 10.5, 10.1, 9.9, 10.3}
	y := []float64{10.0, 9.7, 10.1, 9.8, 9.5, 10.0}

	wantStat := 5.936657514041429
	wantP := 0.0019358364032480116
	wantDF := 5.0
	tol := 1e-12

	res, err := TTestPaired(x, y, PairedTOptions{
		Alternative:     TwoSided,
		ConfidenceLevel: 0.95,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "statistic", res.Statistic, wantStat, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
	requireAlmostEqual(t, "df", res.DF, wantDF, tol)
}

func TestReferenceOneWayANOVA(t *testing.T) {
	groups := [][]float64{
		{8.1, 8.3, 7.9, 8.0, 8.2},
		{8.8, 9.0, 8.7, 8.9, 9.1},
		{7.5, 7.6, 7.4, 7.7, 7.5},
	}

	wantStat := 111.23809523809491
	wantP := 0.00000001796780307526363
	wantDFBetween := 2.0
	wantDFWithin := 12.0
	tol := 1e-9

	res, err := OneWayANOVA(groups...)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "F statistic", res.FStatistic, wantStat, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
	requireAlmostEqual(t, "df_between", float64(res.DFBetween), wantDFBetween, tol)
	requireAlmostEqual(t, "df_within", float64(res.DFWithin), wantDFWithin, tol)
}

func TestReferenceMannWhitney(t *testing.T) {
	x := []float64{14, 15, 16, 17, 18, 19, 20, 21}
	y := []float64{8, 9, 10, 11, 12, 13, 14, 15}

	wantStat := 62.0
	wantP := 0.00160347596869
	tol := 1e-9

	res, err := MannWhitneyU(x, y, MannWhitneyOptions{
		Alternative: TwoSided,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "U statistic", res.U, wantStat, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
}

func TestReferenceWilcoxonSignedRank(t *testing.T) {
	x := []float64{120, 118, 121, 119, 117, 122, 116, 123, 118, 120}
	y := []float64{115, 117, 119, 118, 116, 120, 114, 121, 117, 119}

	wantStat := 55.0
	wantP := 0.004245585155691297
	tol := 1e-9

	res, err := WilcoxonSignedRank(x, y, TwoSided)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "W statistic", res.W, wantStat, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
}

func TestReferenceChiSquare(t *testing.T) {
	table := [][]int{
		{20, 30, 25},
		{22, 28, 35},
	}

	wantStat := 1.2105991822016162
	wantP := 0.5459108521050702
	wantDF := 2.0
	tol := 1e-12

	res, err := ChiSquareIndependence(table) // false = no Yates correction
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "chi-square statistic", res.Statistic, wantStat, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
	requireAlmostEqual(t, "df", float64(res.DF), wantDF, tol)
}

func TestReferenceFisherExact(t *testing.T) {
	table := [2][2]int{
		{8, 2},
		{1, 5},
	}

	wantOddsRatio := 20.0
	wantP := 0.034965034965034975
	tol := 1e-12

	res, err := FishersExact2x2(table, TwoSided)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "odds ratio", res.OddsRatio, wantOddsRatio, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
	requireAlmostEqual(t, "odds ratio", res.OddsRatio, wantOddsRatio, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
}

func TestReferencePAdjustHolm(t *testing.T) {
	p := []float64{0.001, 0.01, 0.03, 0.2, 0.5}

	want := []float64{0.005, 0.04, 0.09, 0.4, 0.5}
	tol := 1e-12

	got, err := AdjustPValues(p, Holm)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireSliceAlmostEqual(t, "adjusted p-values", got, want, tol)
}

func TestReferencePearsonCorrelation(t *testing.T) {
	x := []float64{10, 8, 13, 9, 11, 14, 6, 4, 12, 7}
	y := []float64{8.04, 6.95, 7.58, 8.81, 8.33, 9.96, 7.24, 4.26, 10.84, 4.82}

	wantR := 0.7970815759062527
	wantP := 0.005759692554521912
	wantCILow := 0.3361634481719578
	wantCIHigh := 0.9499584041674188
	tol := 1e-9

	res, err := PearsonCorrelation(x, y, TwoSided, 0.95)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "r", res.R, wantR, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
	requireAlmostEqual(t, "CI low", res.CI.Low, wantCILow, tol)
	requireAlmostEqual(t, "CI high", res.CI.High, wantCIHigh, tol)

	wantPGreater := 0.0028798462772609674
	greater, err := PearsonCorrelation(x, y, Greater, 0.95)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	requireAlmostEqual(t, "p-value (greater)", greater.PValue, wantPGreater, tol)
}

func TestReferenceSpearmanCorrelation(t *testing.T) {
	x := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	y := []float64{2, 1, 4, 3, 7, 5, 6, 9, 8, 10}

	wantRho := 0.9272727272727272
	wantP := 0.00011203450639397587
	wantCILow := 0.7152130503648249
	wantCIHigh := 0.9829930112566714
	tol := 1e-9

	res, err := SpearmanCorrelation(x, y, TwoSided, 0.95)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "rho", res.R, wantRho, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
	requireAlmostEqual(t, "CI low", res.CI.Low, wantCILow, tol)
	requireAlmostEqual(t, "CI high", res.CI.High, wantCIHigh, tol)
}

func TestReferenceKruskalWallis(t *testing.T) {
	g1 := []float64{2.9, 3.0, 2.5, 2.6, 3.2}
	g2 := []float64{3.8, 2.7, 4.0, 2.4, 3.9}
	g3 := []float64{2.8, 3.4, 3.7, 2.2, 2.0}

	wantH := 1.8600000000000065
	wantP := 0.39455371037159986
	tol := 1e-9

	res, err := KruskalWallis(g1, g2, g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "H statistic", res.H, wantH, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
}

func TestReferenceKruskalWallisWithTies(t *testing.T) {
	g1 := []float64{1, 2, 3, 4}
	g2 := []float64{3, 4, 5, 6}
	g3 := []float64{5, 6, 7, 8}

	wantH := 7.645390070921987
	wantP := 0.02186878425495871
	tol := 1e-9

	res, err := KruskalWallis(g1, g2, g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "H statistic", res.H, wantH, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
}

func TestReferenceProportionOneSample(t *testing.T) {
	// Reference values use the null-variance (score-test) convention for the
	// standard error, i.e. se = sqrt(p0*(1-p0)/n) rather than the sample
	// proportion's variance. This matches R's prop.test and the formula in
	// most introductory statistics references.
	wantZ := -0.9999999999999998
	wantP := 0.31731050786291415
	wantPGreater := 0.8413447460685429
	wantPLess := 0.15865525393145707
	wantCILow := 0.3524930229100606
	wantCIHigh := 0.5475069770899395
	tol := 1e-9

	res, err := ProportionOneSample(45, 100, 0.5, ProportionOneSampleOptions{
		Alternative:     TwoSided,
		ConfidenceLevel: 0.95,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	requireAlmostEqual(t, "z", res.Statistic, wantZ, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
	requireAlmostEqual(t, "CI low", res.CI.Low, wantCILow, tol)
	requireAlmostEqual(t, "CI high", res.CI.High, wantCIHigh, tol)

	greater, err := ProportionOneSample(45, 100, 0.5, ProportionOneSampleOptions{
		Alternative:     Greater,
		ConfidenceLevel: 0.95,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	requireAlmostEqual(t, "p-value (greater)", greater.PValue, wantPGreater, tol)

	less, err := ProportionOneSample(45, 100, 0.5, ProportionOneSampleOptions{
		Alternative:     Less,
		ConfidenceLevel: 0.95,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	requireAlmostEqual(t, "p-value (less)", less.PValue, wantPLess, tol)
}

func TestReferenceProportionTwoSample(t *testing.T) {
	wantZ := 1.6427266128719296
	wantP := 0.10043950997963545
	wantDiff := 0.1166666666666667
	wantCILow := -0.021147316981962705
	wantCIHigh := 0.2544806503152961
	tol := 1e-9

	res, err := ProportionTwoSample(45, 100, 30, 90, ProportionTwoSampleOptions{
		Alternative:     TwoSided,
		ConfidenceLevel: 0.95,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	requireAlmostEqual(t, "z", res.Statistic, wantZ, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
	requireAlmostEqual(t, "diff", res.Diff, wantDiff, tol)
	requireAlmostEqual(t, "CI low", res.CI.Low, wantCILow, tol)
	requireAlmostEqual(t, "CI high", res.CI.High, wantCIHigh, tol)
}

func TestReferenceChiSquareGoodnessOfFitUniform(t *testing.T) {
	obs := []float64{18, 22, 16, 24, 20}

	wantChi2 := 2.0
	wantP := 0.7357588823428847
	tol := 1e-9

	res, err := ChiSquareGoodnessOfFit(obs, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	requireAlmostEqual(t, "chi2", res.Statistic, wantChi2, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
}

func TestReferenceChiSquareGoodnessOfFitCustom(t *testing.T) {
	obs := []float64{50, 30, 20}
	exp := []float64{40, 40, 20}

	wantChi2 := 5.0
	wantP := 0.0820849986238988
	tol := 1e-9

	res, err := ChiSquareGoodnessOfFit(obs, exp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	requireAlmostEqual(t, "chi2", res.Statistic, wantChi2, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, tol)
}

func TestReferenceLevenesTest(t *testing.T) {
	g1 := []float64{8.1, 8.3, 7.9, 8.0, 8.2, 8.5}
	g2 := []float64{7.5, 9.0, 6.8, 9.5, 7.0, 8.8}
	g3 := []float64{8.0, 8.1, 7.9, 8.0, 8.1, 7.95}

	wantW := 38.82532159783348
	wantP := 1.1730625046033653e-06
	tol := 1e-7

	res, err := LevenesTest(g1, g2, g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	requireAlmostEqual(t, "W", res.Statistic, wantW, tol)
	requireAlmostEqual(t, "p-value", res.PValue, wantP, 1e-6)
}

func TestReferenceStudentizedRangeCDF(t *testing.T) {
	tol := 1e-4
	cases := []struct {
		q, df float64
		k     int
		want  float64
	}{
		{3.0, 10, 3, 0.8650165848104373},
		{2.5, 20, 4, 0.6827970026274167},
		{4.0, 8, 5, 0.8823515742227231},
		{1.5, 5, 3, 0.42463589484735365},
		{3.5, 30, 6, 0.8361733952381578},
	}
	for _, c := range cases {
		got := studentizedRangeCDF(c.q, c.k, c.df)
		requireAlmostEqual(t, "studentized range CDF", got, c.want, tol)
	}
}

func TestReferenceTukeyHSD(t *testing.T) {
	g1 := []float64{8.1, 8.3, 7.9, 8.0, 8.2}
	g2 := []float64{8.8, 9.0, 8.7, 8.9, 9.1}
	g3 := []float64{7.5, 7.6, 7.4, 7.7, 7.5}

	res, err := TukeyHSD(0.95, g1, g2, g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Comparisons) != 3 {
		t.Fatalf("got %d comparisons, want 3", len(res.Comparisons))
	}

	want := map[[2]int]struct {
		diff, ciLow, ciHigh float64
	}{
		{0, 1}: {-0.8, -1.04451374, -0.55548626},
		{0, 2}: {0.56, 0.31548626, 0.80451374},
		{1, 2}: {1.36, 1.11548626, 1.60451374},
	}

	tol := 1e-5
	for _, c := range res.Comparisons {
		w, ok := want[[2]int{c.GroupI, c.GroupJ}]
		if !ok {
			t.Fatalf("unexpected comparison (%d,%d)", c.GroupI, c.GroupJ)
		}
		requireAlmostEqual(t, "mean diff", c.MeanDiff, w.diff, 1e-12)
		requireAlmostEqual(t, "CI low", c.CI.Low, w.ciLow, tol)
		requireAlmostEqual(t, "CI high", c.CI.High, w.ciHigh, tol)
	}
}

func TestReferenceDunnTest(t *testing.T) {
	g1 := []float64{2.9, 3.0, 2.5, 2.6, 3.2}
	g2 := []float64{3.8, 2.7, 4.0, 2.4, 3.9}
	g3 := []float64{2.8, 3.4, 3.7, 2.2, 2.0}

	res, err := DunnTest(Bonferroni, g1, g2, g3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := map[[2]int]struct{ p, adjP float64 }{
		{0, 1}: {0.288844, 0.866533},
		{0, 2}: {0.832004, 1.0},
		{1, 2}: {0.203092, 0.609275},
	}

	tol := 1e-5
	for _, c := range res.Comparisons {
		w, ok := want[[2]int{c.GroupI, c.GroupJ}]
		if !ok {
			t.Fatalf("unexpected comparison (%d,%d)", c.GroupI, c.GroupJ)
		}
		requireAlmostEqual(t, "p-value", c.PValue, w.p, tol)
		wantAdj := w.adjP
		if wantAdj > 1 {
			wantAdj = 1
		}
		requireAlmostEqual(t, "adjusted p-value", c.AdjustedPValue, wantAdj, tol)
	}
}

func TestReferenceShapiroWilk(t *testing.T) {
	cases := []struct {
		name  string
		x     []float64
		wantW float64
		wantP float64
	}{
		{
			name:  "n5",
			x:     []float64{2.1, 3.4, 1.9, 4.5, 3.0},
			wantW: 0.9415244143556044,
			wantP: 0.6767352294980438,
		},
		{
			name:  "n10",
			x:     []float64{5.1, 4.9, 5.3, 5.0, 4.8, 5.2, 5.4, 4.7, 5.1, 5.0},
			wantW: 0.9839820705366934,
			wantP: 0.9828902222793898,
		},
		{
			name:  "n20_skewed",
			x:     []float64{1, 1, 1, 2, 2, 2, 2, 3, 3, 3, 4, 4, 5, 5, 6, 7, 8, 10, 12, 20},
			wantW: 0.7853402073002983,
			wantP: 0.0005265092668877722,
		},
		{
			name:  "n12_uniform",
			x:     []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			wantW: 0.9668963632914047,
			wantP: 0.8757314433658756,
		},
	}

	for _, c := range cases {
		res, err := ShapiroWilk(c.x)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", c.name, err)
		}
		requireAlmostEqual(t, c.name+" W", res.W, c.wantW, 1e-9)
		requireAlmostEqual(t, c.name+" p-value", res.PValue, c.wantP, 1e-9)
	}
}

func TestReferencePowerTTestTwoSample(t *testing.T) {
	cases := []struct {
		n         int
		d, alpha  float64
		wantPower float64
	}{
		{64, 0.5, 0.05, 0.807430419432557},
		{25, 0.8, 0.05, 0.807430419432557},
	}
	for _, c := range cases {
		got, err := PowerTTestTwoSample(c.n, c.d, c.alpha)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		requireAlmostEqual(t, "power", got, c.wantPower, 1e-9)
	}
}

func TestReferenceSampleSizeTTestTwoSample(t *testing.T) {
	cases := []struct {
		d, alpha, power float64
		wantN           float64
	}{
		{0.5, 0.05, 0.8, 62.79088416571135},
		{0.2, 0.05, 0.9, 525.3709705479278},
		{1.0, 0.01, 0.8, 23.35793629970813},
	}
	for _, c := range cases {
		got, err := SampleSizeTTestTwoSample(c.d, c.alpha, c.power)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		requireAlmostEqual(t, "sample size", got, c.wantN, 1e-4)
	}
}

func TestReferenceCohensH(t *testing.T) {
	got, err := CohensH(0.5, 0.3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	requireAlmostEqual(t, "h", got, 0.4115168460674883, 1e-9)
}

func TestReferencePowerProportionTwoSample(t *testing.T) {
	got, err := PowerProportionTwoSample(50, 0.5, 0.3, 0.05)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	requireAlmostEqual(t, "power", got, 0.5389124796679332, 1e-9)
}

func TestReferenceSampleSizeProportionTwoSample(t *testing.T) {
	got, err := SampleSizeProportionTwoSample(0.5, 0.3, 0.05, 0.8)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	requireAlmostEqual(t, "sample size", got, 92.69608121266178, 1e-3)
}
