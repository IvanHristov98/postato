package number

import (
	"fmt"

	"github.com/IvanHristov98/postato/cluster"
)

const (
	OptimalClusterCount    = 3
	ClusteringRestartCount = 10
)

type gaussianFuzzyNum struct {
	mean   float64
	stdDev float64
}

func SuperClusterToGFNLists(points []*cluster.FuzzyPoint) (map[int][]FuzzyNum, error) {
	gfnLists := make(map[int][]FuzzyNum)

	superCluster := cluster.NewKMeansSuperCluster(points, OptimalClusterCount)
	superCluster.Adjust(ClusteringRestartCount)

	clusteredPoints := superCluster.ClusteredPoints()
	centroids := superCluster.Centroids()
	dimCount, err := superCluster.DimCount()
	if err != nil {
		return nil, fmt.Errorf("Error obtaining dimension count: %s", err)
	}

	for _, centroid := range centroids {
		gfnList := []FuzzyNum{}

		for dim := 0; dim < dimCount; dim++ {
			gfn, err := gfnFromCluster(clusteredPoints, centroid, dim)
			if err != nil {
				return nil, fmt.Errorf("Error obtaining GFN for cluster %d on dim %d: %s", centroid.BestFitClusterIdx, dim, err)
			}

			gfnList = append(gfnList, gfn)
		}

		gfnLists[centroid.BestFitClusterIdx] = gfnList
	}

	return gfnLists, nil
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

func newGaussianFuzzyNum(mean float64, stdDev float64) FuzzyNum {
	return &gaussianFuzzyNum{
		mean:   mean,
		stdDev: stdDev,
	}
}

func (gfn *gaussianFuzzyNum) MembershipDegree(x float64) float64 {
	return 0.0
}

func (gfn *gaussianFuzzyNum) String() string {
	return fmt.Sprintf("mean: %f, std dev: %f", gfn.mean, gfn.stdDev)
}
