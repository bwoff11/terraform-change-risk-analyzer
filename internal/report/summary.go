package report

import (
	"strings"

	"github.com/gwoff/terraform-change-risk-analyzer/internal/risk"
)

// DecideExitCode determines the process exit code based on the highest risk
// level observed and the configured maximum allowed level.
func DecideExitCode(summary risk.Summary, maxAllowedRisk string) int {
	allowed := parseLevel(maxAllowedRisk)
	if summary.MaxLevel > allowed {
		return 1
	}
	return 0
}

func parseLevel(s string) risk.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "none":
		return risk.LevelNone
	case "low":
		return risk.LevelLow
	case "medium":
		return risk.LevelMedium
	case "critical":
		return risk.LevelCritical
	default:
		return risk.LevelHigh
	}
}

