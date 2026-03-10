package tfrisk

// RuleSet is an ordered collection of rules applied to a Change.
type RuleSet struct {
	rules []Rule
}

// NewRuleSet constructs a RuleSet from the given rules.
func NewRuleSet(rules ...Rule) RuleSet {
	return RuleSet{rules: append([]Rule(nil), rules...)}
}

// Score applies the RuleSet to a single Change and produces a ScoreResult.
func (rs RuleSet) Score(c Change) ScoreResult {
	var reasons []string

	action := c.NormalizedAction()
	score, baseReason := actionBaseScore(action)
	if baseReason != "" {
		reasons = append(reasons, baseReason)
	}

	for _, r := range rs.rules {
		if delta, reason, ok := r.Apply(c); ok {
			score += delta
			if reason != "" {
				reasons = append(reasons, reason)
			}
		}
	}

	if score < 0 {
		score = 0
	}

	return ScoreResult{
		Score:    score,
		Severity: SeverityForScore(score),
		Reasons:  reasons,
	}
}

// Engine wraps a RuleSet to provide a clearer public API.
type Engine struct {
	rules RuleSet
}

// NewEngine creates a new Engine with the provided RuleSet.
func NewEngine(rs RuleSet) Engine {
	return Engine{rules: rs}
}

// NewDefaultEngine constructs an Engine using the default rules.
func NewDefaultEngine() Engine {
	return Engine{
		rules: DefaultRuleSet(),
	}
}

// ScoreChange scores a single Change using the engine's rules.
func (e Engine) ScoreChange(c Change) ScoreResult {
	return e.rules.Score(c)
}

// ScoreChanges scores a slice of changes.
func (e Engine) ScoreChanges(changes []Change) []ScoreResult {
	out := make([]ScoreResult, len(changes))
	for i, c := range changes {
		out[i] = e.ScoreChange(c)
	}
	return out
}

// DefaultRuleSet wires up the built-in rules in a fixed order.
func DefaultRuleSet() RuleSet {
	return NewRuleSet(
		dbRule{},
		networkRule{},
		iamRule{},
		lbRule{},
		dnsRule{},
	)
}

