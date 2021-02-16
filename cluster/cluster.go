package cluster

type FuzzySuperCluster interface {
	Adjust(iterCount uint) error
	Clusters() []*Cluster
	BestFitCluster() *Cluster
	SilhouetteCoeff() float64
	ClusteredPoints() []*FuzzyPoint
	Centroids() []*FuzzyPoint
	DimCount() (int, error)
}

type Cluster struct {
	points   []*FuzzyPoint
	centroid *FuzzyPoint
}
