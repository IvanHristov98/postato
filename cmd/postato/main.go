package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	clr "github.com/IvanHristov98/postato/cluster"
	fn "github.com/IvanHristov98/postato/fuzzy/number"
	"github.com/IvanHristov98/postato/plot"
	"github.com/akamensky/argparse"
)

const (
	FeatureCount = 6
	GenDir       = "GENDIR"
	ImageDirName = "image"
)

func main() {
	cfg := parseConfig()

	points, err := parsePoints(cfg.dataset)
	if err != nil {
		log.Fatalf("Error reading points: %s", err)
	}

	fnLists, err := fn.SuperClusterToGFNLists(points)

	if err != nil {
		log.Fatalf("Error building fuzzy numbers from clusters: %s", err)
	}

	cleanUpImages()

	if err := drawAllImages(fnLists); err != nil {
		log.Fatalf("Error drawing fuzzy numbers: %s", err)
	}

	log.Println("Finished")
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

		activity := record[len(record)-1]

		if isNum(activity) {
			activity = ""
		}

		point := clr.NewFuzzyPoint(coords, activity)
		points = append(points, point)
	}

	return points, nil
}

func isNum(val string) bool {
	_, err := strconv.Atoi(val)
	return err == nil
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

func drawAllImages(fnLists map[int][]fn.FuzzyNum) error {
	for clusterIdx, fuzzyNums := range fnLists {
		for fnIdx, fuzzyNum := range fuzzyNums {
			imageName := fmt.Sprintf("fn_%d_%d.png", clusterIdx, fnIdx)
			path, err := imagePath(imageName)
			if err != nil {
				return err
			}

			if err := plot.DrawFuzzyNums(fuzzyNum, -2, 2, path); err != nil {
				return fmt.Errorf("Error drawing fuzzy number %d in cluster %d: %s", fnIdx, clusterIdx, err)
			}
		}
	}

	return nil
}

func cleanUpImages() error {
	dir, err := imageDir()

	if err != nil {
		return fmt.Errorf("Error cleaning up images: %s", err)
	}

	os.RemoveAll(dir)

	return nil
}

func imagePath(name string) (string, error) {
	dir, err := imageDir()

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), name), nil
}

func imageDir() (string, error) {
	dataDir := os.Getenv(GenDir)
	imageDir := fmt.Sprintf("%s%s%s", dataDir, string(os.PathSeparator), ImageDirName)

	if err := os.MkdirAll(imageDir, 0700); err != nil {
		return "", fmt.Errorf("Error creating img dir: %s", err)
	}

	return imageDir, nil
}
