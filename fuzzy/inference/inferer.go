package inference

import "github.com/IvanHristov98/postato/cluster"

type FuzzyInferer interface {
	ClassifyActivity(point *cluster.FuzzyPoint) string
}
