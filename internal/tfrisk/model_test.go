package tfrisk

import "testing"

func TestNormalizedActionReplace(t *testing.T) {
	c := Change{
		Actions: []string{"delete", "create"},
	}
	if got := c.NormalizedAction(); got != ActionReplace {
		t.Fatalf("expected replace, got %q", got)
	}
}

func TestNormalizedActionFirstElement(t *testing.T) {
	c := Change{
		Actions: []string{"update"},
	}
	if got := c.NormalizedAction(); got != ActionUpdate {
		t.Fatalf("expected update, got %q", got)
	}
}

