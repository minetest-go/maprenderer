package maprenderer

import (
	"image/color"
	"math"

	"github.com/fogleman/gg"
)

// limit a number from 0 to 255
func Clamp(num int) uint8 {
	if num < 0 {
		return 0
	}

	if num > 255 {
		return 255
	}

	return uint8(num)
}

// add a color component (darker, lighter)
func AddColorComponent(c *color.RGBA, value int) *color.RGBA {
	return &color.RGBA{
		R: Clamp(int(c.R) + value),
		G: Clamp(int(c.G) + value),
		B: Clamp(int(c.B) + value),
		A: c.A,
	}
}

func AdjustAndFill(dc *gg.Context, r, g, b, adjust int) {
	dc.SetRGB255(
		int(Clamp(r+adjust)),
		int(Clamp(g+adjust)),
		int(Clamp(b+adjust)),
	)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func SortPos(p1, p2 [3]int) ([3]int, [3]int) {
	return [3]int{
			min(p1[0], p2[0]),
			min(p1[1], p2[1]),
			min(p1[2], p2[2]),
		}, [3]int{
			max(p1[0], p2[0]),
			max(p1[1], p2[1]),
			max(p1[2], p2[2]),
		}
}
func NodePosToMapblock(pos [3]int) [3]int {
	// TODO: optimize floating point stuff
	return [3]int{
		int(math.Floor(float64(pos[0]) / 16.0)),
		int(math.Floor(float64(pos[1]) / 16.0)),
		int(math.Floor(float64(pos[2]) / 16.0)),
	}
}

func AddPos(p1, p2 [3]int) [3]int {
	return [3]int{
		p1[0] + p2[0],
		p1[1] + p2[1],
		p1[2] + p2[2],
	}
}
