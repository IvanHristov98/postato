package number

import (
	"fmt"

	"github.com/IvanHristov98/postato/cluster"
)

const (
	// TODO: Tweak values to improve accuracy.
	MinMembershipDegree = 0.15
	BoundWidth          = 0.02
)

type FuzzyNum interface {
	MembershipDegree(x float64) float64
	String() string
}

type FuzzyRule []FuzzyNum

type FuzzyRuleSet map[string]FuzzyRule

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

		if membershipDegree < MinMembershipDegree {
			continue
		}

		// Less than centroid center means min.
		if membershipDegree+BoundWidth > MinMembershipDegree && coord < centroidCoord {
			cumMin += point.Coords[dim]
			minCnt++
		}

		// More than centroid center means max.
		if membershipDegree+BoundWidth > MinMembershipDegree && coord > centroidCoord {
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
