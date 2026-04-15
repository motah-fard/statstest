package statstest

// DefaultConfidence is the default confidence level used in examples and
// recommended as a standard option when the caller does not specify one.
const DefaultConfidence = 0.95

// IsValidConfidence reports whether conf is a valid confidence level in (0, 1).
func IsValidConfidence(conf float64) bool {
	return conf > 0 && conf < 1
}

// MustConfidence returns conf when it is valid, otherwise it returns
// DefaultConfidence.
func MustConfidence(conf float64) float64 {
	if IsValidConfidence(conf) {
		return conf
	}
	return DefaultConfidence
}
