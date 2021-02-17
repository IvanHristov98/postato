package number

import (
	"fmt"
	"math"

	"github.com/IvanHristov98/postato/cluster"
)

type gaussianFuzzyNum struct {
	mean   float64
	stdDev float64
}

func GFNRuleSet(points []*cluster.FuzzyPoint) (FuzzyRuleSet, error) {
	return fuzzyNumRuleSet(points, gfnFromCluster)
}

func (gfn *gaussianFuzzyNum) MembershipDegree(x float64) float64 {
	numer := -math.Pow(x-gfn.mean, 2)
	denom := 2 * math.Pow(gfn.stdDev, 2)
	return math.Exp(numer / denom)
}

func (gfn *gaussianFuzzyNum) String() string {
	return fmt.Sprintf("mean: %f, std dev: %f", gfn.mean, gfn.stdDev)
}

func gfnFromCluster(points []*cluster.FuzzyPoint, centroid *cluster.FuzzyPoint, dim int) (FuzzyNum, error) {
	leftBound, rightBound := clusterBounds(points, centroid, dim)

	mean, err := clusterCenter(leftBound, rightBound)
	if err != nil {
		return nil, fmt.Errorf("Error getting GFN mean: %s", err)
	}

	stdDev, err := clusterWidth(leftBound, rightBound)
	if err != nil {
		return nil, fmt.Errorf("Error obtain GFN standard deviation: %s", err)
	}

	return newGaussianFuzzyNum(mean, stdDev), nil
}

func newGaussianFuzzyNum(mean float64, stdDev float64) FuzzyNum {
	return &gaussianFuzzyNum{
		mean:   mean,
		stdDev: stdDev,
	}
}
