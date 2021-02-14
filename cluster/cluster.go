package cluster

type FuzzySuperCluster interface {
	Adjust(iterCount uint) error
	Clusters() []*Cluster
	BestFitCluster() *Cluster
	SilhouetteCoeff() float64
}

type Cluster struct {
	points   []*FuzzyPoint
	centroid *FuzzyPoint
}
