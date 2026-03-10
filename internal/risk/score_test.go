package risk

import (
	"testing"

	"github.com/gwoff/terraform-change-risk-analyzer/internal/plan"
)

func TestScoreChangesBasic(t *testing.T) {
	pl := plan.Plan{
		ResourceChanges: []plan.ResourceChange{
			{
				Address: "aws_instance.example",
				Mode:    "managed",
				Type:    "aws_instance",
				Change: plan.Change{
					Actions: []string{"create"},
					After: map[string]any{
						"instance_type": "t3.micro",
					},
					AfterUnknown: map[string]any{},
				},
			},
		},
	}

	results, summary := ScoreChanges(pl)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}

	if summary.TotalResources != 1 {
		t.Fatalf("expected TotalResources=1, got %d", summary.TotalResources)
	}

	if results[0].Score.Level == LevelHigh || results[0].Score.Level == LevelCritical {
		t.Fatalf("unexpectedly high risk level for simple create: %s", results[0].Score.Level)
	}
}

