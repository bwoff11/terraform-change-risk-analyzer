package report

import (
	"bytes"
	"strings"
	"testing"
)

func TestBuildSummaryCountsAndOrdering(t *testing.T) {
	risks := []ResourceRisk{
		{Address: "res1", Action: "create", Severity: "LOW", Score: 1},
		{Address: "res2", Action: "update", Severity: "HIGH", Score: 15},
		{Address: "res3", Action: "update", Severity: "MEDIUM", Score: 5},
		{Address: "res4", Action: "delete", Severity: "CRITICAL", Score: 25},
	}

	summary := BuildSummary(risks, 3)

	if summary.Counts.Total != len(risks) {
		t.Fatalf("expected total %d, got %d", len(risks), summary.Counts.Total)
	}
	if summary.Overall.MaxScore != 25 {
		t.Fatalf("expected MaxScore 25, got %d", summary.Overall.MaxScore)
	}
	if summary.Overall.MaxSeverity != "CRITICAL" {
		t.Fatalf("expected MaxSeverity CRITICAL, got %s", summary.Overall.MaxSeverity)
	}
	if len(summary.Top) != 3 {
		t.Fatalf("expected top 3 resources, got %d", len(summary.Top))
	}
	if summary.Top[0].Score < summary.Top[1].Score {
		t.Fatalf("expected top list sorted by score desc, got %v", summary.Top)
	}
}

func TestBuildRecommendationBySeverity(t *testing.T) {
	tests := []struct {
		name   string
		max    string
		total  int
		expect string
	}{
		{"critical", "CRITICAL", 5, "Do NOT proceed"},
		{"high", "HIGH", 5, "manual review"},
		{"medium_small", "MEDIUM", 5, "Medium risk overall"},
		{"low", "LOW", 3, "Risk appears low"},
		{"none", "", 0, "No managed resource changes detected"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Summary{
				Counts: Counts{Total: tt.total},
				Overall: OverallRisk{
					MaxSeverity: tt.max,
				},
			}
			rec := BuildRecommendation(s)
			if !strings.Contains(rec, tt.expect) {
				t.Fatalf("expected recommendation to contain %q, got %q", tt.expect, rec)
			}
		})
	}
}

func TestRenderTextCompact(t *testing.T) {
	summary := Summary{
		Counts: Counts{
			ByAction: map[string]int{"create": 1, "update": 2},
			BySeverity: map[string]int{"LOW": 1, "HIGH": 2},
			Total: 3,
		},
		Overall: OverallRisk{
			MaxSeverity: "HIGH",
			AverageScore: 10.0,
			MaxScore:     20,
		},
		Top: []RankedResource{
			{Address: "aws_db_instance.main", Action: "update", Score: 20, Severity: "HIGH", PrimaryReason: "database change"},
		},
	}
	rec := "Recommendation: test"

	var buf bytes.Buffer
	if err := RenderText(&buf, summary, rec); err != nil {
		t.Fatalf("RenderText returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "tf-risk-report summary") {
		t.Fatalf("expected header in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Actions:") || !strings.Contains(out, "Severities:") {
		t.Fatalf("expected counts section in output, got:\n%s", out)
	}
	if !strings.Contains(out, "SEVERITY") {
		t.Fatalf("expected table header in output, got:\n%s", out)
	}
	if !strings.Contains(out, "Recommendation:") {
		t.Fatalf("expected recommendation section in output, got:\n%s", out)
	}

	lines := strings.Split(out, "\n")
	if len(lines) > 40 {
		t.Fatalf("expected compact output, got %d lines", len(lines))
	}
}

