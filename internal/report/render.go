package report

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/gwoff/terraform-change-risk-analyzer/internal/risk"
)

// RenderText writes a compact, terminal-friendly report to w based on an
// aggregated Summary and recommendation.
func RenderText(w io.Writer, summary Summary, recommendation string) error {
	fmt.Fprintf(w, "tf-risk-report summary (Total: %d, Max severity: %s, Avg score: %.1f)\n",
		summary.Counts.Total,
		summary.Overall.MaxSeverity,
		summary.Overall.AverageScore,
	)

	// Counts section.
	fmt.Fprintln(w, "Actions:", formatCounts(summary.Counts.ByAction))
	fmt.Fprintln(w, "Severities:", formatCounts(summary.Counts.BySeverity))

	// Top risks table.
	if len(summary.Top) == 0 {
		fmt.Fprintln(w, "Top risks: none (no managed resource changes)")
	} else {
		fmt.Fprintln(w, "Top risks:")
		tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
		fmt.Fprintln(tw, "SEVERITY\tSCORE\tACTION\tRESOURCE\tDETAIL")
		for _, r := range summary.Top {
			detail := r.PrimaryReason
			if len(detail) > 80 {
				detail = detail[:77] + "..."
			}
			fmt.Fprintf(tw, "%s\t%d\t%s\t%s\t%s\n",
				r.Severity,
				r.Score,
				r.Action,
				r.Address,
				detail,
			)
		}
		if err := tw.Flush(); err != nil {
			return err
		}
	}

	if recommendation != "" {
		fmt.Fprintln(w)
		fmt.Fprintln(w, recommendation)
	}

	return nil
}

// Render is the legacy entrypoint used by the CLI. It adapts the risk package's
// types into the reporting Summary and delegates to RenderText.
func Render(w io.Writer, changes []risk.ChangeRisk, _ risk.Summary, topN int) error {
	if topN <= 0 {
		topN = len(changes)
	}

	var rr []ResourceRisk
	for _, c := range changes {
		rr = append(rr, ResourceRisk{
			Address:  c.Address,
			Type:     c.ResourceType,
			Action:   c.Action,
			Score:    c.Score.Numeric,
			Severity: strings.ToUpper(c.Score.Level.String()),
			Reasons:  c.Score.Reasons,
		})
	}

	summary := BuildSummary(rr, topN)
	recommendation := BuildRecommendation(summary)
	return RenderText(w, summary, recommendation)
}

func formatCounts(m map[string]int) string {
	if len(m) == 0 {
		return "none"
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var parts []string
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%d", k, m[k]))
	}
	return strings.Join(parts, ", ")
}

