package maprenderer

import "image/color"

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
