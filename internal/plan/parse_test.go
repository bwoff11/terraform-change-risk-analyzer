package plan

import (
	"strings"
	"testing"
)

func TestParseMinimalPlan(t *testing.T) {
	const input = `{
		"resource_changes": [
			{
				"address": "aws_instance.example",
				"mode": "managed",
				"type": "aws_instance",
				"name": "example",
				"provider_name": "registry.terraform.io/hashicorp/aws",
				"change": {
					"actions": ["create"],
					"before": null,
					"after": {
						"instance_type": "t3.micro"
					},
					"after_unknown": {}
				}
			}
		]
	}`

	p, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}

	if len(p.ResourceChanges) != 1 {
		t.Fatalf("expected 1 resource change, got %d", len(p.ResourceChanges))
	}

	rc := p.ResourceChanges[0]
	if rc.Address != "aws_instance.example" {
		t.Errorf("unexpected address: %q", rc.Address)
	}

	if rc.Change.PrimaryAction() != "create" {
		t.Errorf("expected primary action create, got %q", rc.Change.PrimaryAction())
	}
}

