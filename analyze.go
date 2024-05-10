package analyze

import (
	"bufio"
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
	"sort"
)

func CrackFile(file string, topN int) (Colors []uint32, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(bufio.NewReader(f))
	if err != nil {
		return nil, err
	}
	if img == nil {
		return nil, errors.New("nil image object")
	}

	colors := CrackImage(img, topN)
	return colors, nil
}

type Color struct {
	R, G, B uint32
}

type ColorMap struct {
	Count int
	Color Color
}

type ByCount []ColorMap

func (a ByCount) Len() int           { return len(a) }
func (a ByCount) Less(i, j int) bool { return a[i].Count > a[j].Count }
func (a ByCount) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func CrackImage(img image.Image, topN int) (Colors []uint32) {
	colorCounts := make(map[Color]int)

	bounds := img.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			currColor := img.At(x, y)
			r, g, b, _ := currColor.RGBA()
			thisColor := Color{r >> 8, g >> 8, b >> 8}
			colorCounts[thisColor]++
		}
	}

	colorBoard := []ColorMap{}
	for thisColor, count := range colorCounts {
		colorBoard = append(colorBoard, ColorMap{
			Count: count,
			Color: thisColor,
		})
	}

	sort.Sort(ByCount(colorBoard))
	var selectedColors []Color
	var similar bool
	for _, c := range colorBoard {
		curColor := c.Color

		similar = false
		for _, had := range selectedColors {
			if IsSimilarColor(
				color.RGBA{R: uint8(curColor.R), G: uint8(curColor.G), B: uint8(curColor.B)},
				color.RGBA{R: uint8(had.R), G: uint8(had.G), B: uint8(had.B)},
				60) {
				similar = true
				break
			}
		}
		if !similar {
			selectedColors = append(selectedColors, curColor)
		}
	}

	for idx, curColor := range selectedColors {
		if topN > 0 && idx >= topN {
			return
		}
		Colors = append(Colors, uint32(curColor.R)<<16|uint32(curColor.G)<<8|uint32(curColor.B))
	}

	return
}

func IsSimilarColor(c1, c2 color.RGBA, threshold float64) bool {
	rDiff := math.Abs(float64(c1.R) - float64(c2.R))
	gDiff := math.Abs(float64(c1.G) - float64(c2.G))
	bDiff := math.Abs(float64(c1.B) - float64(c2.B))

	colorDiff := math.Sqrt(rDiff*rDiff + gDiff*gDiff + bDiff*bDiff)

	return colorDiff <= threshold
}

func RBA(c uint32) color.Color {
	return color.RGBA{R: uint8((c >> 16) & 0xff), G: uint8((c >> 8) & 0xff), B: uint8(c & 0xff)}
}

func HexColor(c uint32) string {
	return fmt.Sprintf("#%06X", c)
}
