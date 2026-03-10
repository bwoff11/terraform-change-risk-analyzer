package report

import "fmt"

// BuildRecommendation generates a short textual recommendation based on the
// aggregated summary.
func BuildRecommendation(s Summary) string {
	max := s.Overall.MaxSeverity
	total := s.Counts.Total

	switch max {
	case "CRITICAL":
		return "Recommendation: Do NOT proceed with apply. Review all CRITICAL changes before merging."
	case "HIGH":
		return "Recommendation: Proceed only after manual review of HIGH-risk resources listed above."
	case "MEDIUM":
		if total > 20 {
			return fmt.Sprintf("Recommendation: Consider spot-checking MEDIUM changes (total %d), focusing on the highest scores.", total)
		}
		return "Recommendation: Medium risk overall; consider manual review of the highest-scoring resources."
	case "LOW", "", "NONE":
		if total == 0 {
			return "Recommendation: No managed resource changes detected."
		}
		return "Recommendation: Risk appears low for this change set; automated approval is reasonable."
	default:
		return fmt.Sprintf("Recommendation: Overall risk level %s; review the top risks above before proceeding.", max)
	}
}

