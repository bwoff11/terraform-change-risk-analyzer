package plan

// Plan represents the subset of a Terraform plan JSON needed for risk analysis.
type Plan struct {
	ResourceChanges []ResourceChange `json:"resource_changes"`
}

// ResourceChange describes a single resource-level change in the plan.
type ResourceChange struct {
	Address      string `json:"address"`
	Mode         string `json:"mode"`
	Type         string `json:"type"`
	Name         string `json:"name"`
	ProviderName string `json:"provider_name"`
	Change       Change `json:"change"`
}

// Change contains the before/after details for a resource.
type Change struct {
	Actions      []string               `json:"actions"`
	Before       map[string]any         `json:"before"`
	After        map[string]any         `json:"after"`
	AfterUnknown map[string]any         `json:"after_unknown"`
}

// PrimaryAction returns a normalized primary action string derived from the
// Terraform actions array, such as "create", "update", "delete", "replace",
// "no-op", or "read".
func (c Change) PrimaryAction() string {
	if len(c.Actions) == 0 {
		return ""
	}

	// Terraform uses ["delete","create"] to represent a replace.
	if len(c.Actions) == 2 {
		hasDelete := false
		hasCreate := false
		for _, a := range c.Actions {
			switch a {
			case "delete":
				hasDelete = true
			case "create":
				hasCreate = true
			}
		}
		if hasDelete && hasCreate {
			return "replace"
		}
	}

	// In all other cases, default to the first action.
	return c.Actions[0]
}

