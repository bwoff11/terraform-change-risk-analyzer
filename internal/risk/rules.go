package risk

import (
	"fmt"
	"strings"

	"github.com/gwoff/terraform-change-risk-analyzer/internal/plan"
)

type ruleFunc func(rc plan.ResourceChange) (delta int, reason string, matched bool)

var actionBaseScore = map[string]int{
	"no-op":  0,
	"read":   0,
	"create": 20,
	"update": 30,
	"delete": 60,
	"replace": 70,
}

var statefulResourceTypes = map[string]struct{}{
	"aws_db_instance": {},
	"aws_rds_cluster": {},
	"aws_elb":         {},
	"aws_lb":          {},
	"aws_efs_file_system": {},
}

var networkSecurityTypes = map[string]struct{}{
	"aws_security_group":       {},
	"aws_security_group_rule":  {},
	"aws_network_acl":          {},
	"aws_network_acl_rule":     {},
}

var rules = []ruleFunc{
	statefulResourceRule,
	networkExposureRule,
	instanceSizingRule,
	unknownsRule,
}

func statefulResourceRule(rc plan.ResourceChange) (int, string, bool) {
	if _, ok := statefulResourceTypes[rc.Type]; !ok {
		return 0, "", false
	}

	action := rc.Change.PrimaryAction()
	if action == "update" || action == "delete" || action == "replace" {
		return 20, fmt.Sprintf("stateful resource %s %s", rc.Type, action), true
	}

	return 0, "", false
}

func networkExposureRule(rc plan.ResourceChange) (int, string, bool) {
	if _, ok := networkSecurityTypes[rc.Type]; !ok {
		return 0, "", false
	}

	after := rc.Change.After
	before := rc.Change.Before
	if after == nil {
		return 0, "", false
	}

	if exposesWorld(before, after, "cidr_blocks", "0.0.0.0/0") {
		return 30, "network exposure widened to 0.0.0.0/0", true
	}
	if exposesWorld(before, after, "ipv6_cidr_blocks", "::/0") {
		return 30, "network exposure widened to ::/0", true
	}

	return 0, "", false
}

func exposesWorld(before, after map[string]any, key, worldValue string) bool {
	afterVals, ok := after[key]
	if !ok {
		return false
	}

	afterList, ok := afterVals.([]any)
	if !ok {
		return false
	}

	afterHasWorld := containsString(afterList, worldValue)

	var beforeHasWorld bool
	if before != nil {
		if beforeVals, ok := before[key]; ok {
			if beforeList, ok := beforeVals.([]any); ok {
				beforeHasWorld = containsString(beforeList, worldValue)
			}
		}
	}

	return afterHasWorld && !beforeHasWorld
}

func containsString(list []any, target string) bool {
	for _, v := range list {
		if s, ok := v.(string); ok && strings.EqualFold(s, target) {
			return true
		}
	}
	return false
}

func instanceSizingRule(rc plan.ResourceChange) (int, string, bool) {
	after := rc.Change.After
	before := rc.Change.Before
	if after == nil || before == nil {
		return 0, "", false
	}

	keys := []string{"instance_type", "size", "disk_size", "allocated_storage"}
	for _, k := range keys {
		afterVal, okAfter := after[k]
		beforeVal, okBefore := before[k]
		if !okAfter || !okBefore {
			continue
		}
		if fmt.Sprint(afterVal) != fmt.Sprint(beforeVal) {
			return 15, fmt.Sprintf("%s changed (%v -> %v)", k, beforeVal, afterVal), true
		}
	}

	return 0, "", false
}

func unknownsRule(rc plan.ResourceChange) (int, string, bool) {
	if len(rc.Change.AfterUnknown) == 0 {
		return 0, "", false
	}

	return 10, "plan contains unknown values for this resource", true
}

