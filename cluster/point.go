package cluster

import "math"

type FuzzyPoint struct {
	BestFitClusterIdx int
	Coords            []float64
	// Don't need to know how many clusters there are.
	membershipDegrees map[int]float64
}

func NewFuzzyPoint(coords []float64) *FuzzyPoint {
	return &FuzzyPoint{
		BestFitClusterIdx: NoClusterFit,
		Coords:            coords,
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

func (f *FuzzyPoint) hasBestFitCluster() bool {
	return f.BestFitClusterIdx != NoClusterFit
}

func (f *FuzzyPoint) copy(other *FuzzyPoint) {
	f.BestFitClusterIdx = other.BestFitClusterIdx
	f.Coords = other.Coords

	for i := 0; i < len(f.membershipDegrees); i++ {
		f.membershipDegrees[i] = other.membershipDegrees[i]
	}
}
