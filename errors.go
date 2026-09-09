package statstest

import "errors"

var (
	ErrEmptySample            = errors.New("sample must not be empty")
	ErrSampleTooSmall         = errors.New("sample size is too small")
	ErrMismatchedLengths      = errors.New("input slices must have the same length")
	ErrInvalidAlternative     = errors.New("invalid alternative hypothesis")
	ErrInvalidConfidenceLevel = errors.New("confidence level must be in (0, 1)")
	ErrContainsNaN            = errors.New("input contains NaN")
	ErrContainsInf            = errors.New("input contains infinite value")
	ErrInvalidTable           = errors.New("invalid contingency table")
	ErrNegativeCount          = errors.New("counts must be non-negative")
	ErrNoNonZeroDifferences   = errors.New("no non-zero paired differences")
	ErrInvalidPValue          = errors.New("p-values must be in [0, 1]")
	ErrInvalidPAdjustMethod   = errors.New("invalid p-value adjustment method")
	ErrZeroVariance           = errors.New("variance must be greater than zero")
	ErrTooFewGroups           = errors.New("at least two groups are required")
	ErrInsufficientDF         = errors.New("insufficient degrees of freedom")
	ErrInvalidProportion      = errors.New("proportion must be in (0, 1)")
	ErrInvalidCount           = errors.New("count must be between 0 and the number of trials")
	ErrInvalidExpectedCounts  = errors.New("expected counts must be positive")
	ErrSampleTooLarge         = errors.New("sample size is too large")
)
