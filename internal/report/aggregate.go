package report

import "sort"

// BuildSummary aggregates per-resource risks into counts, overall risk and top-N ranking.
func BuildSummary(risks []ResourceRisk, topN int) Summary {
	counts := Counts{
		ByAction:   make(map[string]int),
		BySeverity: make(map[string]int),
	}

	var (
		sumScores      int
		maxScore       int
		maxSeverityStr string
	)

	ranked := make([]RankedResource, 0, len(risks))

	for _, r := range risks {
		counts.Total++
		counts.ByAction[r.Action]++
		counts.BySeverity[r.Severity]++

		sumScores += r.Score
		if r.Score > maxScore {
			maxScore = r.Score
		}

		if maxSeverityStr == "" {
			maxSeverityStr = r.Severity
		} else {
			maxSeverityStr = maxSeverity(maxSeverityStr, r.Severity)
		}

		primary := ""
		if len(r.Reasons) > 0 {
			primary = r.Reasons[0]
			if len(r.Reasons) > 1 {
				primary = primary + " (+" + itoa(len(r.Reasons)-1) + " more)"
			}
		}

		ranked = append(ranked, RankedResource{
			Address:       r.Address,
			Type:          r.Type,
			Action:        r.Action,
			Score:         r.Score,
			Severity:      r.Severity,
			PrimaryReason: primary,
		})
	}

	if len(risks) == 0 {
		return Summary{
			Counts: counts,
			Overall: OverallRisk{
				MaxSeverity: "",
				AverageScore: 0,
				MaxScore:     0,
			},
			Top: nil,
		}
	}

	avg := float64(sumScores) / float64(len(risks))

	sort.Slice(ranked, func(i, j int) bool {
		if ranked[i].Score == ranked[j].Score {
			return ranked[i].Address < ranked[j].Address
		}
		return ranked[i].Score > ranked[j].Score
	})

	if topN > 0 && topN < len(ranked) {
		ranked = ranked[:topN]
	}

	return Summary{
		Counts: counts,
		Overall: OverallRisk{
			MaxSeverity: maxSeverityStr,
			AverageScore: avg,
			MaxScore:     maxScore,
		},
		Top: ranked,
	}
}

// itoa is a tiny integer-to-string helper avoiding fmt for hot paths.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [12]byte
	i := len(buf)
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

