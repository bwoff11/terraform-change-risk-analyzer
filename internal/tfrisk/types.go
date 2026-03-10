package tfrisk

// Severity represents the qualitative severity bucket.
type Severity int

const (
	SeverityLow Severity = iota
	SeverityMedium
	SeverityHigh
	SeverityCritical
)

func (s Severity) String() string {
	switch s {
	case SeverityMedium:
		return "MEDIUM"
	case SeverityHigh:
		return "HIGH"
	case SeverityCritical:
		return "CRITICAL"
	default:
		return "LOW"
	}
}

// ScoreResult is the outcome of scoring a single Change.
type ScoreResult struct {
	Score    int
	Severity Severity
	Reasons  []string
}

// Severity thresholds for numeric scores.
const (
	lowMax      = 2
	mediumMax   = 9
	highMax     = 19
	criticalMin = 20
)

// SeverityForScore maps a numeric score to a Severity bucket.
func SeverityForScore(score int) Severity {
	switch {
	case score >= criticalMin:
		return SeverityCritical
	case score > highMax:
		return SeverityHigh
	case score > mediumMax:
		return SeverityHigh
	case score > lowMax:
		return SeverityMedium
	default:
		return SeverityLow
	}
}

