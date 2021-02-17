package number

import (
	"fmt"

	"github.com/IvanHristov98/postato/cluster"
)

func NewFuzzyRuleSet(fuzzyNumType string, points []*cluster.FuzzyPoint) (FuzzyRuleSet, error) {
	switch fuzzyNumType {
	case GaussianFuzzyNum:
		return GFNRuleSet(points)
	case TriangularFuzzyNum:
		return TFNRuleSet(points)
	default:
		return nil, fmt.Errorf("Invalid fuzzy num type provided %s", fuzzyNumType)
	}
}
