package maprenderer

import (
	"image/color"

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

func ClampInt(num int) int {
	if num < 0 {
		return 0
	}

	if num > 255 {
		return 255
	}

	return num
}

func AdjustColorComponent(c uint8, adj int) int {
	num := int(c) + adj
	return ClampInt(num)
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

func SortPos(p1, p2 *Pos) (*Pos, *Pos) {
	return &Pos{
			min(p1[0], p2[0]),
			min(p1[1], p2[1]),
			min(p1[2], p2[2]),
		}, &Pos{
			max(p1[0], p2[0]),
			max(p1[1], p2[1]),
			max(p1[2], p2[2]),
		}
}
