package risk

// Level represents the qualitative risk classification.
type Level int

const (
	LevelNone Level = iota
	LevelLow
	LevelMedium
	LevelHigh
	LevelCritical
)

func (l Level) String() string {
	switch l {
	case LevelLow:
		return "low"
	case LevelMedium:
		return "medium"
	case LevelHigh:
		return "high"
	case LevelCritical:
		return "critical"
	default:
		return "none"
	}
}

// Score is the numeric and qualitative risk score for a single change.
type Score struct {
	Level   Level
	Numeric int
	Reasons []string
}

// ChangeRisk describes the scored risk of a single resource change.
type ChangeRisk struct {
	Address      string
	Action       string
	ResourceType string
	Score        Score
}

// Summary aggregates risk across all resource changes.
type Summary struct {
	TotalResources int
	ByLevel        map[Level]int
	MaxLevel       Level
}

