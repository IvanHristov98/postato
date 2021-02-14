package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	clr "github.com/IvanHristov98/postato/cluster"
	"github.com/akamensky/argparse"
)

const FeatureCount = 6

func main() {
	cfg := parseConfig()

	points, err := parsePoints(cfg.dataset)
	if err != nil {
		log.Fatalf("Error reading points: %s", err)
	}

	// Some of the positions could have multiple clusters.
	clusterCnt := 2
	c := clr.NewKMeansSuperCluster(points, clusterCnt)
	c.Adjust(10)
}

type config struct {
	dataset string
}

func parseConfig() *config {
	parser := argparse.NewParser("postato", "Guesses human body position")

	d := parser.String("d", "dataset", &argparse.Options{Required: true, Help: "Path to training dataset. Must be a CSV."})

	if err := parser.Parse(os.Args); err != nil {
		log.Fatalf("Error parsing arguments: %s", err)
	}

	return &config{
		dataset: *d,
	}
}

func parsePoints(path string) ([]*clr.FuzzyPoint, error) {
	points := []*clr.FuzzyPoint{}

	records, err := readCSVFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error reading points: %s", err)
	}

	for i, record := range records {
		coords := []float64{}

		for j := 0; j < FeatureCount; j++ {
			col := record[j]
			coord, err := strconv.ParseFloat(col, 64)
			if err != nil {
				return nil, fmt.Errorf("Error reading value %s in record %d: %s", col, i, err)
			}

			coords = append(coords, coord)
		}

		point := clr.NewFuzzyPoint(coords)
		points = append(points, point)
	}

	return points, nil
}

func readCSVFile(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to read input file: %s", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("Error parsing CSV records: %s", err)
	}

	return records, nil
}
