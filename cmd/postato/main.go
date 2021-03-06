package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	clr "github.com/IvanHristov98/postato/cluster"
	"github.com/IvanHristov98/postato/fuzzy/inference"
	"github.com/IvanHristov98/postato/fuzzy/number"
	fn "github.com/IvanHristov98/postato/fuzzy/number"
	"github.com/IvanHristov98/postato/plot"
	"github.com/akamensky/argparse"
)

const (
	FeatureCount   = 6
	GridBound      = 2
	GenDir         = "GENDIR"
	ImageDirName   = "image"
	FoldCrossCount = 10
)

type config struct {
	dataset string
	fnType  string
}

func main() {
	rand.Seed(time.Now().UnixNano())

	parser := argparse.NewParser("postato", "Guesses human body position")

	d := parser.String("d", "dataset", &argparse.Options{Required: true, Help: "Path to training dataset. Must be a CSV."})

	fuzzyNumTypes := []string{number.GaussianFuzzyNum, number.TriangularFuzzyNum}
	t := parser.Selector("t", "type", fuzzyNumTypes, &argparse.Options{Required: false, Default: number.GaussianFuzzyNum})

	drawCmd := parser.NewCommand("draw", "Draws fuzzy numbers.")
	testCmd := parser.NewCommand("test", "Runs a fold cross validation of the fuzzy inference system.")

	if err := parser.Parse(os.Args); err != nil {
		log.Fatalf("Error parsing arguments: %s", err)
	}

	cfg := &config{dataset: *d, fnType: *t}

	if drawCmd.Happened() {
		drawFuzzyNumbers(cfg)
	} else if testCmd.Happened() {
		crossFold(cfg)
	}
}

func drawFuzzyNumbers(cfg *config) {
	points, err := parsePoints(cfg.dataset)
	if err != nil {
		log.Fatalf("Error reading points: %s", err)
	}

	fuzzyRuleSet, err := fn.NewFuzzyRuleSet(cfg.fnType, points)

	if err != nil {
		log.Fatalf("Error building fuzzy numbers from clusters: %s", err)
	}

	cleanUpImages()

	if err := drawAllImages(fuzzyRuleSet); err != nil {
		log.Fatalf("Error drawing fuzzy numbers: %s", err)
	}

	log.Println("Fuzzy number drawing completed.")
}

func crossFold(cfg *config) {
	points, err := parsePoints(cfg.dataset)
	if err != nil {
		log.Fatalf("Error reading points: %s", err)
	}

	rand.Shuffle(len(points), func(i, j int) { points[i], points[j] = points[j], points[i] })
	cumAccuracy := 0.0

	for i := 0; i < FoldCrossCount; i++ {
		trainingPoints := append(points[:int(len(points)*i/FoldCrossCount)], points[int(len(points)*(i+1)/FoldCrossCount):]...)

		fuzzyRuleSet, err := fn.NewFuzzyRuleSet(cfg.fnType, trainingPoints)

		if err != nil {
			log.Fatalf("Error building fuzzy numbers from clusters: %s", err)
		}

		inferer := inference.NewMamdaniInferer(fuzzyRuleSet)

		testPoints := points[int(len(points)*i/10):int(len(points)*(i+1)/FoldCrossCount)]

		successCnt := 0

		for _, testPoint := range testPoints {
			activity := inferer.ClassifyActivity(testPoint)

			if testPoint.Activity == activity {
				successCnt++
			}
		}

		accuracy := 100.0 * float64(successCnt) / float64(len(testPoints))
		log.Printf("Accuracy of fold cross %d iteration is %2.f perc.\n", i, accuracy)

		cumAccuracy += accuracy
	}

	avgAccuracy := cumAccuracy / float64(FoldCrossCount)
	log.Printf("Average accuracy of %d fold cross is %2.f perc.\n", FoldCrossCount, avgAccuracy)
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

func drawAllImages(fuzzyRuleSet fn.FuzzyRuleSet) error {
	for activity, fuzzyNums := range fuzzyRuleSet {
		for fnIdx, fuzzyNum := range fuzzyNums {
			imageName := fmt.Sprintf("fn_%s_%d.png", activity, fnIdx)
			path, err := imagePath(imageName)
			if err != nil {
				return err
			}

			if err := plot.DrawFuzzyNums(fuzzyNum, -GridBound, GridBound, fnIdx, activity, path); err != nil {
				return fmt.Errorf("Error drawing fuzzy number %d in activity %s: %s", fnIdx, activity, err)
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
