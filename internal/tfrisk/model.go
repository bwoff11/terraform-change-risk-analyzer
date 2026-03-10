package tfrisk

// Action represents a normalized Terraform change action.
type Action string

const (
	ActionCreate  Action = "create"
	ActionUpdate  Action = "update"
	ActionDelete  Action = "delete"
	ActionReplace Action = "replace"
	ActionRead    Action = "read"
	ActionNoop    Action = "no-op"
)

// Change describes a single Terraform resource change for scoring.
type Change struct {
	Address string
	Type    string
	Actions []string

	// Metadata can hold optional fields (e.g. provider, tags) for future rules.
	Metadata map[string]any
}

// NormalizedAction maps the raw Terraform actions array to a single Action.
// A combination of delete and create is treated as replace, otherwise the
// first element is used as the primary action.
func (c Change) NormalizedAction() Action {
	if len(c.Actions) == 0 {
		return ActionNoop
	}

	// Terraform expresses replace as both delete and create.
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
		return ActionReplace
	}

	switch c.Actions[0] {
	case "create":
		return ActionCreate
	case "update":
		return ActionUpdate
	case "delete":
		return ActionDelete
	case "read":
		return ActionRead
	case "no-op":
		return ActionNoop
	default:
		return Action(c.Actions[0])
	}
}

