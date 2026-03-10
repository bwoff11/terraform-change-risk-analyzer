package plan

import (
	"encoding/json"
	"fmt"
	"io"
)

// Parse reads a Terraform plan JSON document from r and decodes the subset of
// fields required for risk analysis.
func Parse(r io.Reader) (Plan, error) {
	var p Plan

	dec := json.NewDecoder(r)
	if err := dec.Decode(&p); err != nil {
		return Plan{}, fmt.Errorf("decode plan JSON: %w", err)
	}

	if p.ResourceChanges == nil {
		return Plan{}, fmt.Errorf("input does not appear to be a Terraform plan JSON (missing resource_changes)")
	}

	return p, nil
}

