package integration

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gwoff/terraform-change-risk-analyzer/internal/plan"
	"github.com/gwoff/terraform-change-risk-analyzer/internal/report"
	"github.com/gwoff/terraform-change-risk-analyzer/internal/risk"
)

func TestEndToEndSimplePlan(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "plan_simple.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	pl, err := plan.Parse(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	changes, summary := risk.ScoreChanges(pl)
	if summary.TotalResources != 1 {
		t.Fatalf("expected 1 resource, got %d", summary.TotalResources)
	}

	var buf bytes.Buffer
	if err := report.Render(&buf, changes, summary, 10); err != nil {
		t.Fatalf("Render returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "aws_instance.example") {
		t.Fatalf("report output did not contain expected resource address, got:\n%s", out)
	}
}

func TestEndToEndHighRiskPlan(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "plan_high_risk.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read fixture: %v", err)
	}

	pl, err := plan.Parse(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	_, summary := risk.ScoreChanges(pl)
	if summary.MaxLevel < risk.LevelHigh {
		t.Fatalf("expected at least high risk, got %s", summary.MaxLevel.String())
	}

	code := report.DecideExitCode(summary, "medium")
	if code == 0 {
		t.Fatalf("expected non-zero exit code for high risk plan with max-allowed-risk=medium")
	}
}

