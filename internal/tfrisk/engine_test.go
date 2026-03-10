package tfrisk

import "testing"

func TestEngineScoreChange_ActionAndTypeRules(t *testing.T) {
	engine := NewDefaultEngine()

	tests := []struct {
		name     string
		change   Change
		minScore int
		wantType string
	}{
		{
			name: "simple create low risk",
			change: Change{
				Address: "aws_instance.example",
				Type:    "aws_instance",
				Actions: []string{"create"},
			},
			minScore: 1,
			wantType: "aws_instance",
		},
		{
			name: "db update high risk",
			change: Change{
				Address: "aws_db_instance.main",
				Type:    "aws_db_instance",
				Actions: []string{"update"},
			},
			minScore: 3 + 20, // base update + DB bump
			wantType: "aws_db_instance",
		},
		{
			name: "network security group",
			change: Change{
				Address: "aws_security_group.web",
				Type:    "aws_security_group",
				Actions: []string{"update"},
			},
			minScore: 3 + 15, // base update + network bump
			wantType: "aws_security_group",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := engine.ScoreChange(tt.change)
			if res.Score < tt.minScore {
				t.Fatalf("expected score at least %d, got %d", tt.minScore, res.Score)
			}
			if res.Severity == SeverityLow && res.Score >= 10 {
				t.Fatalf("expected severity to increase for high scores, got LOW for score %d", res.Score)
			}
			if len(res.Reasons) == 0 {
				t.Fatalf("expected at least one reason, got none")
			}
		})
	}
}

