package main

import (
	"fmt"

	clr "github.com/IvanHristov98/postato/cluster"
)

func main() {
	p := clr.NewFuzzyPoint([]float64{2.0})
	clusterCnt := 1

	c := clr.NewKMeansSuperCluster([]*clr.FuzzyPoint{p}, clusterCnt)
	c.Adjust(10)

	fmt.Print(p.BestFitClusterIdx)
}
