package main

import (
	"fmt"

	clr "github.com/IvanHristov98/postato/cluster"
)

func main() {
	points := []*clr.FuzzyPoint{
		clr.NewFuzzyPoint([]float64{2.0}),
		clr.NewFuzzyPoint([]float64{2.0}),
		clr.NewFuzzyPoint([]float64{5.0}),
		clr.NewFuzzyPoint([]float64{4.0}),
		clr.NewFuzzyPoint([]float64{2.0}),
	}
	clusterCnt := 2

	c := clr.NewKMeansSuperCluster(points, clusterCnt)
	c.Adjust(10)

	for _, point := range points {
		fmt.Println(point.BestFitClusterIdx, point.MembershipDegree(0), point.MembershipDegree(1))
	}
}
