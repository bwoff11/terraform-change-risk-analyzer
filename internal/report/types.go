package report

// ResourceRisk is a reporting-friendly view of a single resource's risk.
type ResourceRisk struct {
	Address  string
	Type     string
	Action   string
	Score    int
	Severity string
	Reasons  []string
}

// Counts holds aggregate counts by action and severity.
type Counts struct {
	ByAction   map[string]int
	BySeverity map[string]int
	Total      int
}

// OverallRisk summarizes the overall risk across all resources.
type OverallRisk struct {
	MaxSeverity string
	AverageScore float64
	MaxScore     int
}

// RankedResource is a compact view of a resource for top-N output.
type RankedResource struct {
	Address       string
	Type          string
	Action        string
	Score         int
	Severity      string
	PrimaryReason string
}

// Summary aggregates counts, overall risk, and the top risky resources.
type Summary struct {
	Counts  Counts
	Overall OverallRisk
	Top     []RankedResource
}

var severityOrder = map[string]int{
	"LOW":      0,
	"MEDIUM":   1,
	"HIGH":     2,
	"CRITICAL": 3,
}

func maxSeverity(a, b string) string {
	ra, oka := severityOrder[a]
	rb, okb := severityOrder[b]
	if !oka && !okb {
		return a
	}
	if !oka {
		return b
	}
	if !okb {
		return a
	}
	if rb > ra {
		return b
	}
	return a
}

