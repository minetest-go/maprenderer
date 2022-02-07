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
