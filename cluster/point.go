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

func (p *FuzzyPoint) Clone() *FuzzyPoint {
	cloneMembershipDegrees := make(map[int]float64, len(p.membershipDegrees))

	for k, v := range p.membershipDegrees {
		cloneMembershipDegrees[k] = v
	}

	return &FuzzyPoint{
		BestFitClusterIdx: p.BestFitClusterIdx,
		Coords:            p.Coords,
		membershipDegrees: cloneMembershipDegrees,
	}
}

func (p *FuzzyPoint) Dist(other *FuzzyPoint) float64 {
	dist := 0.0

	for dim, coord := range p.Coords {
		otherCoord := other.Coords[dim]

		dist += math.Pow(coord-otherCoord, 2)
	}

	return math.Sqrt(dist)
}

func (p *FuzzyPoint) DimCount() int {
	return len(p.Coords)
}

func (p *FuzzyPoint) hasBestFitCluster() bool {
	return p.BestFitClusterIdx != NoClusterFit
}
