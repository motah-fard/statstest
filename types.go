package statstest

// Alternative represents the alternative hypothesis for a statistical test.
type Alternative string

const (
	TwoSided Alternative = "two-sided"
	Less     Alternative = "less"
	Greater  Alternative = "greater"
)

// ConfidenceInterval represents a confidence interval at a given level.
type ConfidenceInterval struct {
	Level float64
	Low   float64
	High  float64
}

// PAdjustMethod identifies a multiple-testing correction method.
type PAdjustMethod string

const (
	Bonferroni        PAdjustMethod = "bonferroni"
	Holm              PAdjustMethod = "holm"
	BenjaminiHochberg PAdjustMethod = "benjamini-hochberg"
)

// TwoSampleTOptions configures a two-sample t-test.
type TwoSampleTOptions struct {
	EqualVariance   bool
	Alternative     Alternative
	ConfidenceLevel float64
}

// PairedTOptions configures a paired t-test.
type PairedTOptions struct {
	Alternative     Alternative
	ConfidenceLevel float64
}

// MannWhitneyOptions configures a Mann-Whitney U test.
type MannWhitneyOptions struct {
	Alternative             Alternative
	UseContinuityCorrection bool
}

// TTestOneSampleResult contains the result of a one-sample t-test.
type TTestOneSampleResult struct {
	Statistic   float64
	PValue      float64
	DF          float64
	Mean        float64
	NullMean    float64
	MeanDiff    float64
	CI          ConfidenceInterval
	Method      string
	Alternative Alternative
}

// TTestTwoSampleResult contains the result of a two-sample t-test.
type TTestTwoSampleResult struct {
	Statistic   float64
	PValue      float64
	DF          float64
	Mean1       float64
	Mean2       float64
	MeanDiff    float64
	CI          ConfidenceInterval
	Method      string
	Alternative Alternative
}

// TTestPairedResult contains the result of a paired t-test.
type TTestPairedResult struct {
	Statistic   float64
	PValue      float64
	DF          float64
	MeanDiff    float64
	CI          ConfidenceInterval
	Method      string
	Alternative Alternative
}

// MannWhitneyResult contains the result of a Mann-Whitney U test.
type MannWhitneyResult struct {
	U           float64
	PValue      float64
	Z           float64
	Method      string
	Alternative Alternative
}

// WilcoxonSignedRankResult contains the result of a Wilcoxon signed-rank test.
type WilcoxonSignedRankResult struct {
	W           float64
	PValue      float64
	Z           float64
	Method      string
	Alternative Alternative
}

// ChiSquareResult contains the result of a chi-square test.
type ChiSquareResult struct {
	Statistic float64
	PValue    float64
	DF        int
	Expected  [][]float64
	Method    string
}

// FishersExactResult contains the result of Fisher's exact test for a 2x2 table.
type FishersExactResult struct {
	OddsRatio   float64
	PValue      float64
	Method      string
	Alternative Alternative
}

// OneWayANOVAResult contains the result of a one-way ANOVA.
type OneWayANOVAResult struct {
	FStatistic float64
	PValue     float64
	DFBetween  int
	DFWithin   int
	SSTotal    float64
	SSBetween  float64
	SSWithin   float64
	EtaSquared float64
	Method     string
}

// KruskalWallisResult contains the result of a Kruskal-Wallis H test.
type KruskalWallisResult struct {
	H      float64
	PValue float64
	DF     int
	Method string
}

// PairwiseComparison contains one pairwise group comparison from a
// post-hoc test such as TukeyHSD. GroupI and GroupJ are indices into the
// groups passed to the test.
type PairwiseComparison struct {
	GroupI    int
	GroupJ    int
	MeanDiff  float64
	Statistic float64
	PValue    float64
	CI        ConfidenceInterval
}

// TukeyHSDResult contains the result of Tukey's Honestly Significant
// Difference post-hoc test.
type TukeyHSDResult struct {
	Comparisons []PairwiseComparison
	Method      string
}

// DunnComparison contains one pairwise group comparison from Dunn's test.
// GroupI and GroupJ are indices into the groups passed to the test.
type DunnComparison struct {
	GroupI         int
	GroupJ         int
	Statistic      float64
	PValue         float64
	AdjustedPValue float64
}

// DunnTestResult contains the result of Dunn's post-hoc test.
type DunnTestResult struct {
	Comparisons  []DunnComparison
	AdjustMethod PAdjustMethod
	Method       string
}

// ChiSquareGoodnessOfFitResult contains the result of a chi-square
// goodness-of-fit test.
type ChiSquareGoodnessOfFitResult struct {
	Statistic float64
	PValue    float64
	DF        int
	Expected  []float64
	Method    string
}

// LeveneResult contains the result of Levene's test for equality of
// variances.
type LeveneResult struct {
	Statistic float64
	PValue    float64
	DFBetween int
	DFWithin  int
	Method    string
}

// ShapiroWilkResult contains the result of the Shapiro-Wilk normality
// test.
type ShapiroWilkResult struct {
	W      float64
	PValue float64
	Method string
}

// CorrelationResult contains the result of a correlation test.
type CorrelationResult struct {
	R           float64
	PValue      float64
	DF          float64
	N           int
	CI          ConfidenceInterval
	Method      string
	Alternative Alternative
}

// ProportionOneSampleOptions configures a one-sample proportion test.
type ProportionOneSampleOptions struct {
	Alternative             Alternative
	ConfidenceLevel         float64
	UseContinuityCorrection bool
}

// ProportionTwoSampleOptions configures a two-sample proportion test.
type ProportionTwoSampleOptions struct {
	Alternative             Alternative
	ConfidenceLevel         float64
	UseContinuityCorrection bool
}

// ProportionResult contains the result of a one-sample proportion test.
type ProportionResult struct {
	Proportion  float64
	NullValue   float64
	Statistic   float64
	PValue      float64
	CI          ConfidenceInterval
	Method      string
	Alternative Alternative
}

// ProportionTwoSampleResult contains the result of a two-sample proportion test.
type ProportionTwoSampleResult struct {
	Proportion1 float64
	Proportion2 float64
	Diff        float64
	Statistic   float64
	PValue      float64
	CI          ConfidenceInterval
	Method      string
	Alternative Alternative
}
