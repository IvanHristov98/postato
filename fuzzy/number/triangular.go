package number

import (
	"fmt"

	"github.com/IvanHristov98/postato/cluster"
)

type triangularFuzzyNum struct {
	left   float64
	center float64
	right  float64
}

func TFNRuleSet(points []*cluster.FuzzyPoint) (FuzzyRuleSet, error) {
	return fuzzyNumRuleSet(points, tfnFromCluster)
}

func (t *triangularFuzzyNum) MembershipDegree(x float64) float64 {
	switch {
	case x < t.left:
		return 0.0
	case x < t.center:
		return (x - t.left) / (t.center - t.left)
	case x < t.right:
		return (t.right - x) / (t.right - t.center)
	default:
		return 0.0
	}
}

func (t *triangularFuzzyNum) String() string {
	return fmt.Sprintf("left: %2.f, center: %2.f, right: %2.f", t.left, t.center, t.right)
}

func tfnFromCluster(points []*cluster.FuzzyPoint, centroid *cluster.FuzzyPoint, dim int) (FuzzyNum, error) {
	leftBound, rightBound := clusterBounds(points, centroid, dim)

	mean, err := clusterCenter(leftBound, rightBound)
	if err != nil {
		return nil, fmt.Errorf("Error getting GFN mean: %s", err)
	}

	return newTriangularFuzzyNum(leftBound, mean, rightBound), nil
}

func newTriangularFuzzyNum(left, center, right float64) FuzzyNum {
	return &triangularFuzzyNum{
		left:   left,
		center: center,
		right:  right,
	}
}
