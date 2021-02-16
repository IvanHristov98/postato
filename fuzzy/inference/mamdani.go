package inference

import (
	"github.com/IvanHristov98/postato/cluster"
	"github.com/IvanHristov98/postato/fuzzy/number"
)

const (
	MinMembershipDegree = 0.0
	MaxMembershipDegree = 1.0
)

type mamdaniInferer struct {
	ruleSet number.FuzzyRuleSet
}

func NewMamdaniInferer(ruleSet number.FuzzyRuleSet) FuzzyInferer {
	return &mamdaniInferer{
		ruleSet: ruleSet,
	}
}

func (m *mamdaniInferer) ClassifyActivity(point *cluster.FuzzyPoint) string {
	maxMembershipDegree := MinMembershipDegree
	bestFitActivity := ""

	for activity := range m.ruleSet {
		membershipDegree := m.activityMembershipDegree(point, activity)

		if maxMembershipDegree < membershipDegree {
			maxMembershipDegree = membershipDegree
			bestFitActivity = activity
		}
	}

	return bestFitActivity
}

func (m *mamdaniInferer) activityMembershipDegree(point *cluster.FuzzyPoint, activity string) float64 {
	rule := m.ruleSet[activity]
	minMembershipDegree := MaxMembershipDegree

	for i, fuzzyNum := range rule {
		membershipDegree := fuzzyNum.MembershipDegree(point.Coords[i])

		if membershipDegree < minMembershipDegree {
			minMembershipDegree = membershipDegree
		}
	}

	return minMembershipDegree
}
