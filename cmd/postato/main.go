package main

import (
	"encoding/csv"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
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

	c := clr.NewKMeansSuperCluster(points, 3)
	c.Adjust(4)
	cnt := 0
	counts := map[string]int{"sitting": 0.0, "lying": 0.0, "standing": 0.0}

	for _, point := range points {
		if point.BestFitClusterIdx != 2 {
			continue
		}

		cnt++
		counts[point.Activity]++
	}

	for key, val := range counts {
		fmt.Printf("%s: %v, ", key, float64(val)/float64(cnt))
	}

	fmt.Println("Size", cnt)

	drawMembershipDegrees(points)
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

func drawMembershipDegrees(points []*clr.FuzzyPoint) {
	md := "md.png"
	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	green := color.RGBA{255, 255, 255, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{green}, image.ZP, draw.Src)

	rect := image.Rect(127, 127, 128, 128)
	red := color.RGBA{255, 0, 0, 255}

	draw.Draw(img, rect, &image.Uniform{red}, image.ZP, draw.Src)

	file, err := os.Create(md)
	if err != nil {
		panic(err)
	}

	png.Encode(file, img)
}
