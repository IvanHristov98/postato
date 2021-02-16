package cluster

type FuzzySuperCluster interface {
	Adjust(iterCount uint) error
	SilhouetteCoeff() float64
	ClusteredPoints() []*FuzzyPoint
	Centroids() []*FuzzyPoint
	DimCount() (int, error)
}
