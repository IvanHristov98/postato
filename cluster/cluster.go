package cluster

const NoClusterFit = -1

type FuzzySuperCluster interface {
	Adjust(iterCount uint) error
	Clusters() []*Cluster
	BestFitCluster() *Cluster
}

type Cluster struct {
	points   []*FuzzyPoint
	centroid *FuzzyPoint
}
