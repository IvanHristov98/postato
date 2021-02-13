package cluster

import (
	"fmt"
	"math"
	"math/rand"
)

const (
	InitialCentroidProb         = 1.0
	MaxCentroidSamplingAttempts = 10
)

type kMeansSuperCluster struct {
	points         []*FuzzyPoint
	clusterCount   int
	minClusterDist float64
}

func NewKMeansSuperCluster(points []*FuzzyPoint, clusterCount int) FuzzySuperCluster {
	return &kMeansSuperCluster{
		points:         points,
		clusterCount:   clusterCount,
		minClusterDist: math.Inf(0),
	}
}

func (k *kMeansSuperCluster) Adjust(iterCount uint) error {
	for i := 0; i < int(iterCount); i++ {
		clonedPoints := k.clonePoints()
		k.clusterize(clonedPoints)

		// TODO: Intercluster dist metric and improve points.
	}

	return nil
}

func (k *kMeansSuperCluster) Clusters() []*Cluster {
	return nil
}

func (k *kMeansSuperCluster) BestFitCluster() *Cluster {
	return &Cluster{}
}

func (k *kMeansSuperCluster) clonePoints() []*FuzzyPoint {
	clonedPoints := []*FuzzyPoint{}

	for _, point := range k.points {
		clonedPoint := point.Clone()
		clonedPoints = append(clonedPoints, clonedPoint)
	}

	return clonedPoints
}

func (k *kMeansSuperCluster) clusterize(points []*FuzzyPoint) error {
	centroids, err := k.initialCentroids(points)

	if err != nil {
		return fmt.Errorf("Error building cluster: %s", err.Error())
	}

	clusterSizes := []int{}

	for i := 0; i < k.clusterCount; i++ {
		centroids[i].BestFitClusterIdx = i
		clusterSizes = append(clusterSizes, 0)
	}

	madeAdjustments := true

	for madeAdjustments {
		madeAdjustments = k.adjustClusters(points, centroids, clusterSizes)

		if err := k.adjustCentroids(points, centroids, clusterSizes); err != nil {
			return fmt.Errorf("Eror building cluster: %s", err.Error())
		}
	}

	return nil
}

func (k *kMeansSuperCluster) initialCentroids(points []*FuzzyPoint) ([]*FuzzyPoint, error) {
	centroids := []*FuzzyPoint{}

	probabilities := []float64{}
	probSum := 0.0

	for range points {
		probabilities = append(probabilities, InitialCentroidProb)
		probSum += InitialCentroidProb
	}

	for i := 0; i < k.clusterCount; i++ {
		index, err := randDistributionIndex(probabilities, probSum, 0)
		if err != nil {
			return nil, fmt.Errorf("Error selecting centroid %d: %s", i, err.Error())
		}

		centroid := points[index].Clone()
		centroids = append(centroids, centroid)

		probSum = 0.0

		for j, point := range points {
			_, minDist := bestFitCluster(centroids, point)

			probabilities[j] = math.Pow(minDist, 2)
			probSum += probabilities[j]
		}
	}

	return centroids, nil
}

func (k *kMeansSuperCluster) adjustClusters(points, centroids []*FuzzyPoint, clusterSizes []int) (madeAdjustments bool) {
	for _, point := range points {
		bestFitClusterIdx, _ := bestFitCluster(centroids, point)

		if bestFitClusterIdx != point.BestFitClusterIdx {
			clusterSizes[bestFitClusterIdx]++

			if point.hasBestFitCluster() {
				clusterSizes[point.BestFitClusterIdx]--
			}

			point.BestFitClusterIdx = bestFitClusterIdx
			madeAdjustments = true
		}

	}

	return madeAdjustments
}

func (k *kMeansSuperCluster) adjustCentroids(points, centroids []*FuzzyPoint, clusterSizes []int) error {
	for i, centroid := range centroids {
		if clusterSizes[i] == 0 {
			continue
		}

		center, err := k.clusterCenter(points, i)
		if err != nil {
			return fmt.Errorf("Error adjusting cluster centers: %s", err.Error())
		}

		centroid.Coords = center
	}

	return nil
}

func (k *kMeansSuperCluster) clusterCenter(points []*FuzzyPoint, clusterIdx int) ([]float64, error) {
	dimCount, err := k.dimCount()
	if err != nil {
		return nil, fmt.Errorf("Error finding cluster center: %s", err.Error())
	}

	overallCoords := make([]float64, dimCount)
	// TODO: Possibly remove since they should be initialized with a size of zero.
	for i := 0; i < dimCount; i++ {
		overallCoords[i] = 0
	}

	pointsCount := 0

	for _, point := range points {
		if point.BestFitClusterIdx != clusterIdx {
			continue
		}

		for i, coord := range point.Coords {
			overallCoords[i] += coord
		}

		pointsCount++
	}

	if pointsCount == 0 {
		for i := 0; i < dimCount; i++ {
			overallCoords[i] = math.Inf(0)
		}

		return overallCoords, nil
	}

	for i := 0; i < dimCount; i++ {
		overallCoords[i] /= float64(pointsCount)
	}

	return overallCoords, nil
}

func (k *kMeansSuperCluster) dimCount() (int, error) {
	if len(k.points) == 0 {
		return 0, fmt.Errorf("No points to clusterize")
	}
	point := k.points[0]
	return point.DimCount(), nil
}

func bestFitCluster(centroids []*FuzzyPoint, point *FuzzyPoint) (int, float64) {
	bestFitClusterIdx := point.BestFitClusterIdx
	minDist := math.Inf(0)

	for clusterIdx, centroid := range centroids {
		dist := point.Dist(centroid)

		if dist < minDist {
			minDist = dist
			bestFitClusterIdx = clusterIdx
		}

		// TODO: Set membership degree logic.
	}

	return bestFitClusterIdx, minDist
}

func randDistributionIndex(probabilities []float64, probSum float64, attempts int) (int, error) {
	if attempts > MaxCentroidSamplingAttempts {
		return -1, fmt.Errorf("Random distribution index not selected")
	}

	randNum := rand.Float64() * probSum
	sum := 0.0

	for i := 0; i < len(probabilities); i++ {
		sum += probabilities[i]

		if sum > randNum {
			return i, nil
		}
	}

	return randDistributionIndex(probabilities, probSum, attempts+1)
}
