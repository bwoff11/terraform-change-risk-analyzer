package report

import (
	"testing"

	"github.com/gwoff/terraform-change-risk-analyzer/internal/risk"
)

func TestDecideExitCode(t *testing.T) {
	summary := risk.Summary{MaxLevel: risk.LevelHigh}

	if code := DecideExitCode(summary, "medium"); code != 1 {
		t.Fatalf("expected non-zero exit when allowed=medium and max=high, got %d", code)
	}

	if code := DecideExitCode(summary, "high"); code != 0 {
		t.Fatalf("expected zero exit when allowed=high and max=high, got %d", code)
	}
}

