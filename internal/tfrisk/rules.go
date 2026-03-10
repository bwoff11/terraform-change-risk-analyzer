package tfrisk

// Rule applies a scoring adjustment for a given Change.
type Rule interface {
	Apply(change Change) (delta int, reason string, ok bool)
}

// actionBaseScore maps normalized actions to their base scores.
func actionBaseScore(a Action) (int, string) {
	switch a {
	case ActionCreate:
		return 1, "base risk from action create"
	case ActionUpdate:
		return 3, "base risk from action update"
	case ActionDelete:
		return 8, "base risk from action delete"
	case ActionReplace:
		return 10, "base risk from action replace"
	default:
		return 0, "no base risk for action " + string(a)
	}
}

// dbRule increases risk for database-related resources.
type dbRule struct{}

func (r dbRule) Apply(c Change) (int, string, bool) {
	if _, ok := dbTypes[c.Type]; !ok {
		return 0, "", false
	}
	return 20, "database-related resource type", true
}

// networkRule increases risk for network/security-related resources.
type networkRule struct{}

func (r networkRule) Apply(c Change) (int, string, bool) {
	if _, ok := networkTypes[c.Type]; !ok {
		return 0, "", false
	}
	return 15, "network/security group or firewall resource", true
}

// iamRule increases risk for IAM-related resources.
type iamRule struct{}

func (r iamRule) Apply(c Change) (int, string, bool) {
	if _, ok := iamTypes[c.Type]; !ok {
		return 0, "", false
	}
	return 12, "IAM roles or policies", true
}

// lbRule increases risk for load balancer resources.
type lbRule struct{}

func (r lbRule) Apply(c Change) (int, string, bool) {
	if _, ok := lbTypes[c.Type]; !ok {
		return 0, "", false
	}
	return 10, "load balancer resource", true
}

// dnsRule increases risk for DNS-related resources.
type dnsRule struct{}

func (r dnsRule) Apply(c Change) (int, string, bool) {
	if _, ok := dnsTypes[c.Type]; !ok {
		return 0, "", false
	}
	return 8, "DNS or routing resource", true
}

var dbTypes = map[string]struct{}{
	"aws_db_instance":             {},
	"aws_rds_cluster":             {},
	"google_sql_database_instance": {},
}

var networkTypes = map[string]struct{}{
	"aws_security_group":      {},
	"aws_security_group_rule": {},
	"aws_network_acl":         {},
}

var iamTypes = map[string]struct{}{
	"aws_iam_role":                  {},
	"aws_iam_policy":                {},
	"aws_iam_role_policy":           {},
	"aws_iam_role_policy_attachment": {},
}

var lbTypes = map[string]struct{}{
	"aws_elb": {},
	"aws_lb":  {},
}

var dnsTypes = map[string]struct{}{
	"aws_route53_record":        {},
	"google_dns_record_set":     {},
}

