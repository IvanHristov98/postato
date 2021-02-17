package number

import (
	"fmt"

	"github.com/IvanHristov98/postato/cluster"
)

const (
	// TODO: Tweak values to improve accuracy.
	MinViableMembershipDegree = 0.05
	BoundWidth                = 0.02
	OptimalClusterCount       = 3
	ClusteringRestartCount    = 10
	GaussianFuzzyNum          = "gaussian"
	TriangularFuzzyNum        = "triangular"
)

type FuzzyNum interface {
	MembershipDegree(x float64) float64
	String() string
}

type FuzzyRule []FuzzyNum

type FuzzyRuleSet map[string]FuzzyRule

type superClusterToFNConverter func(points []*cluster.FuzzyPoint, centroid *cluster.FuzzyPoint, dim int) (FuzzyNum, error)

func fuzzyNumRuleSet(points []*cluster.FuzzyPoint, converter superClusterToFNConverter) (FuzzyRuleSet, error) {
	ruleSet := make(FuzzyRuleSet)

	superCluster := cluster.NewKMeansSuperCluster(points, OptimalClusterCount)
	superCluster.Adjust(ClusteringRestartCount)

	clusteredPoints := superCluster.ClusteredPoints()
	centroids := superCluster.Centroids()
	dimCount, err := superCluster.DimCount()
	if err != nil {
		return nil, fmt.Errorf("Error obtaining dimension count: %s", err)
	}

	for _, centroid := range centroids {
		rule := FuzzyRule{}

		for dim := 0; dim < dimCount; dim++ {
			gfn, err := converter(clusteredPoints, centroid, dim)
			if err != nil {
				return nil, fmt.Errorf("Error obtaining GFN for cluster %d on dim %d: %s", centroid.BestFitClusterIdx, dim, err)
			}

			rule = append(rule, gfn)
		}

		ruleSet[centroid.Activity] = rule
	}

	return ruleSet, nil
}

func clusterBounds(points []*cluster.FuzzyPoint, centroid *cluster.FuzzyPoint, dim int) (float64, float64) {
	cumMin := 0.0
	minCnt := 0
	cumMax := 0.0
	maxCnt := 0

	centroidCoord := centroid.Coords[dim]

	for _, point := range points {
		// The best fit cluster index of a centroid should always be the cluster it belongs to.
		membershipDegree := point.MembershipDegree(centroid.BestFitClusterIdx)
		coord := point.Coords[dim]

		if membershipDegree < MinViableMembershipDegree {
			continue
		}

		// Less than centroid center means min.
		if membershipDegree+BoundWidth > MinViableMembershipDegree && coord < centroidCoord {
			cumMin += point.Coords[dim]
			minCnt++
		}

		// More than centroid center means max.
		if membershipDegree+BoundWidth > MinViableMembershipDegree && coord > centroidCoord {
			cumMax += point.Coords[dim]
			maxCnt++
		}
	}

	if minCnt == 0 {
		fmt.Printf("Min for clr %d dim %d NaN\n", centroid.BestFitClusterIdx, dim)
	}

	if maxCnt == 0 {
		fmt.Printf("Max for clr %d dim %d NaN\n", centroid.BestFitClusterIdx, dim)
	}

	min := cumMin / float64(minCnt)
	max := cumMax / float64(maxCnt)

	return min, max
}

func clusterCenter(leftBound float64, rightBound float64) (float64, error) {
	if rightBound <= leftBound {
		return 0.0, fmt.Errorf("Left bound greater than or equal to right bound")
	}

	return (rightBound + leftBound) / 2, nil
}

func clusterWidth(leftBound float64, rightBound float64) (float64, error) {
	if rightBound <= leftBound {
		return 0.0, fmt.Errorf("Left bound greater than or equal to right bound")
	}

	return rightBound - leftBound, nil
}
