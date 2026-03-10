package report

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"

	"github.com/gwoff/terraform-change-risk-analyzer/internal/risk"
)

// Render writes a human-readable risk report to w.
func Render(w io.Writer, changes []risk.ChangeRisk, summary risk.Summary, topN int) error {
	if topN <= 0 {
		topN = len(changes)
	}

	fmt.Fprintf(w, "tf-risk-report\n\n")

	fmt.Fprintf(w, "Summary:\n")
	fmt.Fprintf(w, "  Total managed resources: %d\n", summary.TotalResources)
	fmt.Fprintf(w, "  Highest risk level: %s\n", summary.MaxLevel.String())
	fmt.Fprintf(w, "  Counts by level:\n")
	for _, lvl := range []risk.Level{risk.LevelCritical, risk.LevelHigh, risk.LevelMedium, risk.LevelLow, risk.LevelNone} {
		count := summary.ByLevel[lvl]
		fmt.Fprintf(w, "    %-8s: %d\n", lvl.String(), count)
	}
	fmt.Fprintln(w)

	if len(changes) == 0 {
		fmt.Fprintln(w, "No managed resource changes detected.")
		return nil
	}

	sort.Slice(changes, func(i, j int) bool {
		if changes[i].Score.Numeric == changes[j].Score.Numeric {
			return changes[i].Address < changes[j].Address
		}
		return changes[i].Score.Numeric > changes[j].Score.Numeric
	})

	if topN > len(changes) {
		topN = len(changes)
	}

	fmt.Fprintf(w, "Top %d highest-risk changes:\n", topN)

	tw := tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "LEVEL\tSCORE\tACTION\tRESOURCE\tDETAILS")

	for i := 0; i < topN; i++ {
		cr := changes[i]

		var detail string
		if len(cr.Score.Reasons) > 0 {
			detail = cr.Score.Reasons[0]
			if len(cr.Score.Reasons) > 1 {
				remaining := len(cr.Score.Reasons) - 1
				detail = fmt.Sprintf("%s (+%d more)", detail, remaining)
			}
		}

		fmt.Fprintf(tw, "%s\t%d\t%s\t%s\t%s\n",
			cr.Score.Level.String(),
			cr.Score.Numeric,
			cr.Action,
			cr.Address,
			detail,
		)
	}

	if err := tw.Flush(); err != nil {
		return err
	}

	if topN < len(changes) {
		fmt.Fprintf(w, "\nOnly showing top %d of %d changes.\n", topN, len(changes))
	}

	return nil
}

