package tfrisk

import (
	"strings"
	"testing"
)

func TestEngine_ScoreChange_TableDriven(t *testing.T) {
	eng := NewDefaultEngine()

	tests := []struct {
		name               string
		change             Change
		wantMinScore       int
		wantSeverity       Severity
		wantReasonContains string
	}{
		{
			name: "create instance low risk",
			change: Change{
				Address: "aws_instance.example",
				Type:    "aws_instance",
				Actions: []string{"create"},
			},
			wantMinScore:       1,
			wantSeverity:       SeverityLow,
			wantReasonContains: "action create",
		},
		{
			name: "update instance low risk",
			change: Change{
				Address: "aws_instance.example",
				Type:    "aws_instance",
				Actions: []string{"update"},
			},
			wantMinScore:       3,
			wantSeverity:       SeverityMedium,
			wantReasonContains: "action update",
		},
		{
			name: "delete instance higher risk",
			change: Change{
				Address: "aws_instance.example",
				Type:    "aws_instance",
				Actions: []string{"delete"},
			},
			wantMinScore:       8,
			wantSeverity:       SeverityMedium,
			wantReasonContains: "action delete",
		},
		{
			name: "replace instance",
			change: Change{
				Address: "aws_instance.example",
				Type:    "aws_instance",
				Actions: []string{"delete", "create"},
			},
			wantMinScore:       10,
			wantSeverity:       SeverityHigh,
			wantReasonContains: "action replace",
		},
		{
			name: "database replacement",
			change: Change{
				Address: "aws_db_instance.main",
				Type:    "aws_db_instance",
				Actions: []string{"delete", "create"},
			},
			wantMinScore:       10 + 20,
			wantSeverity:       SeverityCritical,
			wantReasonContains: "database",
		},
		{
			name: "IAM update",
			change: Change{
				Address: "aws_iam_role.app",
				Type:    "aws_iam_role",
				Actions: []string{"update"},
			},
			wantMinScore:       3 + 12,
			wantSeverity:       SeverityHigh,
			wantReasonContains: "IAM",
		},
		{
			name: "security group update",
			change: Change{
				Address: "aws_security_group.web",
				Type:    "aws_security_group",
				Actions: []string{"update"},
			},
			wantMinScore:       3 + 15,
			wantSeverity:       SeverityHigh,
			wantReasonContains: "network/security",
		},
		{
			name: "no-op change",
			change: Change{
				Address: "aws_instance.example",
				Type:    "aws_instance",
				Actions: []string{"no-op"},
			},
			wantMinScore:       0,
			wantSeverity:       SeverityLow,
			wantReasonContains: "",
		},
		{
			name: "empty actions",
			change: Change{
				Address: "aws_instance.example",
				Type:    "aws_instance",
				Actions: []string{},
			},
			wantMinScore:       0,
			wantSeverity:       SeverityLow,
			wantReasonContains: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := eng.ScoreChange(tt.change)

			if res.Score < tt.wantMinScore {
				t.Fatalf("score=%d, want at least %d", res.Score, tt.wantMinScore)
			}
			if res.Severity != tt.wantSeverity {
				t.Fatalf("severity=%v, want %v", res.Severity, tt.wantSeverity)
			}
			if tt.wantReasonContains != "" && !reasonContains(res.Reasons, tt.wantReasonContains) {
				t.Fatalf("reasons=%v, expected to contain %q", res.Reasons, tt.wantReasonContains)
			}
		})
	}
}

func reasonContains(reasons []string, sub string) bool {
	for _, r := range reasons {
		if strings.Contains(r, sub) {
			return true
		}
	}
	return false
}

