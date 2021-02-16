package plot

import (
	"fmt"
	"math/rand"
	"os"

	fn "github.com/IvanHristov98/postato/fuzzy/number"
	"github.com/fogleman/gg"
)

const (
	ImgWidth            = 1200
	ImgHeight           = 600
	HOffset             = 50.0
	VOffset             = 50.0
	LineWidth           = 2.0
	LineAlpha           = 1.0
	GridLineWidth       = 1.0
	GridLineAlpha       = 0.2
	CurveAlpha          = 0.75
	CurveWidth          = 3.0
	TextOffset          = 10
	GridCount           = 10
	DataPointCount      = 500
	MaxMembershipDegree = 1.0
	MinMembershipDegree = 0.0

	Font        = "PTSans-Regular.ttf"
	DataDir     = "DATADIR"
	FontDirName = "font"
)

type dataPoint struct {
	x float64
	y float64
}

func DrawFuzzyNums(num fn.FuzzyNum, low, high float64, imagePath string) error {
	dc, err := drawingCanvas(low, high)
	if err != nil {
		return fmt.Errorf("Error initializing canvas: %s", err)
	}

	drawFuzzyNum(dc, num, low, high)

	if err := dc.SavePNG(imagePath); err != nil {
		return fmt.Errorf("Error saving png %s: %s", imagePath, err)
	}

	return nil
}

func drawingCanvas(low, high float64) (*gg.Context, error) {
	dc := gg.NewContext(ImgWidth, ImgHeight)

	dc.SetRGBA(1, 1, 1, 1)
	dc.Clear()

	drawGrid(dc, low, high)

	return dc, nil
}

func drawGrid(dc *gg.Context, low, high float64) error {
	if err := loadFont(dc); err != nil {
		return fmt.Errorf("Error loading font: %s", err)
	}

	drawGridBounds(dc)
	drawHorizontalSubgrid(dc, low, high)

	return nil
}

func drawGridBounds(dc *gg.Context) {
	dc.SetRGBA(0, 0, 0, LineAlpha)

	drawLine(dc, LineWidth, HOffset, VOffset, float64(ImgWidth)-HOffset, VOffset)
	drawLine(dc, LineWidth, HOffset, VOffset, HOffset, float64(ImgHeight)-VOffset)
	drawLine(dc, LineWidth, HOffset, float64(ImgHeight)-VOffset, float64(ImgWidth)-HOffset, float64(ImgHeight)-VOffset)
	drawLine(dc, LineWidth, float64(ImgWidth)-HOffset, VOffset, float64(ImgWidth)-HOffset, float64(ImgHeight)-VOffset)
}

func drawHorizontalSubgrid(dc *gg.Context, low, high float64) {
	width := high - low
	subgridWidth := width / float64(GridCount)

	widthInPixels := gridWidth()
	gridWidthInPixels := widthInPixels / float64(GridCount)

	dc.SetRGBA(0, 0, 0, GridLineAlpha)

	for i := 0; float64(i)*subgridWidth <= width; i++ {
		xInPixels := gridWidthInPixels*float64(i) + HOffset
		drawLine(dc, GridLineWidth, xInPixels, VOffset, xInPixels, float64(ImgHeight)-VOffset)

		x := float64(i)*subgridWidth + low
		label := fmt.Sprintf("%.2f", x)
		dc.SetRGBA(0, 0, 0, 1)
		dc.DrawStringAnchored(label, xInPixels, float64(ImgHeight)-VOffset+TextOffset, 0.5, 0.5)
	}
}

func loadFont(dc *gg.Context) error {
	font := fontLocation()
	if err := dc.LoadFontFace(font, 12); err != nil {
		return fmt.Errorf("Error loading font: %s", err)
	}

	return nil
}

func fontLocation() string {
	dataDir := os.Getenv(DataDir)
	fontDir := fmt.Sprintf("%s%s%s", dataDir, string(os.PathSeparator), FontDirName)

	return fmt.Sprintf("%s%s%s", fontDir, string(os.PathSeparator), Font)
}

func drawFuzzyNum(dc *gg.Context, num fn.FuzzyNum, low, high float64) error {
	dpDelta := (high - low) / DataPointCount
	dataPoints := []*dataPoint{}

	for i := 0; i < DataPointCount; i++ {
		x := float64(i)*dpDelta + low
		y := num.MembershipDegree(x)

		dp := &dataPoint{x: x, y: y}
		dataPoints = append(dataPoints, dp)
	}

	if err := drawCurve(dc, dataPoints, low, high); err != nil {
		return fmt.Errorf("Error drawing point %s", num)
	}

	return nil
}

func drawCurve(dc *gg.Context, dataPoints []*dataPoint, low, high float64) error {
	if len(dataPoints) < 2 {
		return fmt.Errorf("A curve consists of at least 2 points")
	}

	dc.SetRGBA(rand.Float64(), rand.Float64(), rand.Float64(), CurveAlpha)

	for i := 0; i < len(dataPoints)-1; i++ {
		currDP := dataPoints[i]
		nextDP := dataPoints[i+1]

		currPosX := pointPosX(currDP, low, high)
		nextPosX := pointPosX(nextDP, low, high)

		currPosY := pointPosY(currDP)
		nextPosY := pointPosY(nextDP)

		drawLine(dc, CurveWidth, currPosX, currPosY, nextPosX, nextPosY)
		dc.Fill()
	}

	return nil
}

func pointPosX(dp *dataPoint, low, high float64) float64 {
	widthInPixels := gridWidth()
	relPos := (dp.x - low) / (high - low)

	return widthInPixels*relPos + HOffset
}

func pointPosY(dp *dataPoint) float64 {
	heightInPixels := gridHeight()
	relPos := (dp.y - MinMembershipDegree) / (MaxMembershipDegree - MinMembershipDegree)

	return heightInPixels + VOffset - heightInPixels*relPos
}

func gridWidth() float64 {
	return float64(ImgWidth) - 2*HOffset
}

func gridHeight() float64 {
	return float64(ImgHeight) - 2*VOffset
}

func drawLine(dc *gg.Context, width, x1, y1, x2, y2 float64) {
	dc.SetLineWidth(width)
	dc.DrawLine(x1, y1, x2, y2)
	dc.Stroke()

	dc.Fill()
}
