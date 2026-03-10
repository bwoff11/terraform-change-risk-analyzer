package tfrisk

import "testing"

func TestSeverityForScoreThresholds(t *testing.T) {
	tests := []struct {
		score int
		want  Severity
	}{
		{0, SeverityLow},
		{2, SeverityLow},
		{3, SeverityMedium},
		{9, SeverityMedium},
		{10, SeverityHigh},
		{19, SeverityHigh},
		{20, SeverityCritical},
		{30, SeverityCritical},
	}

	for _, tt := range tests {
		if got := SeverityForScore(tt.score); got != tt.want {
			t.Errorf("SeverityForScore(%d) = %v, want %v", tt.score, got, tt.want)
		}
	}
}

