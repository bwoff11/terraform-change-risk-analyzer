package risk

import (
	"github.com/gwoff/terraform-change-risk-analyzer/internal/plan"
)

// ScoreChanges evaluates all resource changes in a plan and returns per-change
// risk assessments along with an aggregate summary.
func ScoreChanges(p plan.Plan) ([]ChangeRisk, Summary) {
	results := make([]ChangeRisk, 0, len(p.ResourceChanges))
	summary := Summary{
		ByLevel: make(map[Level]int),
	}

	for _, rc := range p.ResourceChanges {
		if rc.Mode != "managed" {
			continue
		}

		action := rc.Change.PrimaryAction()
		base := actionBaseScore[action]

		score := Score{
			Numeric: base,
		}
		if base > 0 {
			score.Reasons = append(score.Reasons, "base risk from action "+action)
		}

		for _, rule := range rules {
			if delta, reason, matched := rule(rc); matched {
				score.Numeric += delta
				score.Reasons = append(score.Reasons, reason)
			}
		}

		if score.Numeric < 0 {
			score.Numeric = 0
		}
		if score.Numeric > 100 {
			score.Numeric = 100
		}

		score.Level = levelForScore(score.Numeric)

		cr := ChangeRisk{
			Address:      rc.Address,
			Action:       action,
			ResourceType: rc.Type,
			Score:        score,
		}
		results = append(results, cr)

		summary.TotalResources++
		summary.ByLevel[score.Level]++
		if score.Level > summary.MaxLevel {
			summary.MaxLevel = score.Level
		}
	}

	return results, summary
}

func levelForScore(n int) Level {
	switch {
	case n == 0:
		return LevelNone
	case n <= 25:
		return LevelLow
	case n <= 50:
		return LevelMedium
	case n <= 75:
		return LevelHigh
	default:
		return LevelCritical
	}
}

