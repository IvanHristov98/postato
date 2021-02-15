package cluster

import "math"

const NoCluster = -1

type FuzzyPoint struct {
	BestFitClusterIdx int
	Coords            []float64
	Activity          string
	// Don't need to know how many clusters there are.
	membershipDegrees map[int]float64
}

func NewFuzzyPoint(coords []float64, activity string) *FuzzyPoint {
	return &FuzzyPoint{
		BestFitClusterIdx: NoCluster,
		Coords:            coords,
		Activity:          activity,
		membershipDegrees: make(map[int]float64),
	}
}

func (f *FuzzyPoint) Clone() *FuzzyPoint {
	cloneMembershipDegrees := make(map[int]float64, len(f.membershipDegrees))

	for k, v := range f.membershipDegrees {
		cloneMembershipDegrees[k] = v
	}

	return &FuzzyPoint{
		BestFitClusterIdx: f.BestFitClusterIdx,
		Coords:            f.Coords,
		Activity:          f.Activity,
		membershipDegrees: cloneMembershipDegrees,
	}
}

func (f *FuzzyPoint) Dist(other *FuzzyPoint) float64 {
	dist := 0.0

	for dim, coord := range f.Coords {
		otherCoord := other.Coords[dim]

		dist += math.Pow(coord-otherCoord, 2)
	}

	return math.Sqrt(dist)
}

func (f *FuzzyPoint) DimCount() int {
	return len(f.Coords)
}

func (f *FuzzyPoint) MembershipDegree(clusterIdx int) float64 {
	return f.membershipDegrees[clusterIdx]
}

func (f *FuzzyPoint) setMembershipDegree(centroids []*FuzzyPoint) {
	totalMembership := 0.0

	for _, centroid := range centroids {
		dist := f.Dist(centroid)
		if dist == 0 {
			f.membershipDegrees[centroid.BestFitClusterIdx] = 1
		} else {
			f.membershipDegrees[centroid.BestFitClusterIdx] = 1 / (dist * dist)
		}

		totalMembership += f.membershipDegrees[centroid.BestFitClusterIdx]
	}

	for _, centroid := range centroids {
		f.membershipDegrees[centroid.BestFitClusterIdx] /= totalMembership
	}
}

func (f *FuzzyPoint) hasBestFitCluster() bool {
	return f.BestFitClusterIdx != NoCluster
}

func (f *FuzzyPoint) copy(other *FuzzyPoint) {
	f.BestFitClusterIdx = other.BestFitClusterIdx
	f.Coords = other.Coords
	f.Activity = other.Activity

	for i := 0; i < len(f.membershipDegrees); i++ {
		f.membershipDegrees[i] = other.membershipDegrees[i]
	}
}

func (f *FuzzyPoint) nearestClusterIdx() int {
	nearestClr := NoCluster
	maxMembershipDegree := 0.0

	for clusterIdx, membershipDegree := range f.membershipDegrees {
		if clusterIdx == f.BestFitClusterIdx {
			continue
		}

		if maxMembershipDegree < membershipDegree {
			nearestClr = clusterIdx
			maxMembershipDegree = membershipDegree
		}
	}

	return nearestClr
}
