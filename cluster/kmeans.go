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
		// Cloning points to keep original ones intact
		clonedPoints := k.clonePoints()
		centroids, err := k.clusterize(clonedPoints)

		if err != nil {
			return fmt.Errorf("Error adjusting super cluster: %s", err.Error())
		}

		dist := overallClusterDist(centroids, clonedPoints)

		if dist < k.minClusterDist {
			//log.Printf("Encounetered a better cluster with dist %f", dist)

			k.minClusterDist = dist
			k.copyPoints(clonedPoints)

			for _, points := range k.points {
				points.setMembershipDegree(centroids)
			}
		}
	}

	return nil
}

func (k *kMeansSuperCluster) Clusters() []*Cluster {
	return nil
}

func (k *kMeansSuperCluster) BestFitCluster() *Cluster {
	return &Cluster{}
}

func (k *kMeansSuperCluster) SilhouetteCoeff() float64 {
	if k.clusterCount == 1 {
		return 0.0
	}

	cumSilhouetteCoeff := 0.0

	for _, point := range k.points {
		intraCumDist := 0.0
		neighbourCumDist := 0.0
		intraCnt := 0
		neighbourCnt := 0

		nearestNeighbour := point.nearestClusterIdx()

		for _, otherPoint := range k.points {
			if point == otherPoint {
				continue
			}

			if otherPoint.BestFitClusterIdx == point.BestFitClusterIdx {
				dist := point.Dist(otherPoint)
				intraCumDist += dist
				intraCnt++
			} else if otherPoint.BestFitClusterIdx == nearestNeighbour {
				dist := point.Dist(otherPoint)
				neighbourCumDist += dist
				neighbourCnt++
			}
		}

		intraDistMean := intraCumDist / float64(intraCnt)
		neighbourDistMean := neighbourCumDist / float64(neighbourCnt)

		cumSilhouetteCoeff += (neighbourDistMean - intraDistMean) / math.Max(neighbourDistMean, intraDistMean)
	}

	return cumSilhouetteCoeff / float64(len(k.points))
}

func (k *kMeansSuperCluster) clonePoints() []*FuzzyPoint {
	clonedPoints := []*FuzzyPoint{}

	for _, point := range k.points {
		clonedPoint := point.Clone()
		clonedPoints = append(clonedPoints, clonedPoint)
	}

	return clonedPoints
}

func (k *kMeansSuperCluster) clusterize(points []*FuzzyPoint) ([]*FuzzyPoint, error) {
	centroids, err := k.initialCentroids(points)

	if err != nil {
		return nil, fmt.Errorf("Error building cluster: %s", err.Error())
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
			return nil, fmt.Errorf("Eror building cluster: %s", err.Error())
		}
	}

	return centroids, nil
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
		centroid.BestFitClusterIdx = i

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

func (k *kMeansSuperCluster) copyPoints(otherPoints []*FuzzyPoint) {
	for i, point := range k.points {
		point.copy(otherPoints[i])
	}
}

func overallClusterDist(centroids, points []*FuzzyPoint) float64 {
	overallDist := 0.0

	for _, centroid := range centroids {
		overallDist += clusterDist(centroid, points)
	}

	return overallDist
}

func clusterDist(centroid *FuzzyPoint, points []*FuzzyPoint) float64 {
	dist := 0.0

	for _, point := range points {
		if point.BestFitClusterIdx == centroid.BestFitClusterIdx {
			dist += point.Dist(centroid)
		}
	}

	return dist
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
