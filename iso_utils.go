package maprenderer

import (
	"fmt"
	"image"
	"image/color"

	"github.com/minetest-go/types"
)

func GetIsometricImageSize(size *types.Pos, cube_len int) (int, int) {
	width := (size.X() * cube_len / 2) +
		(size.Z() * cube_len / 2)

	height := (size.X() * cube_len / 4) +
		(size.Y() * cube_len / 2) +
		(size.Z() * cube_len / 4)

	return width, height
}

func GetIsoCenterCubeOffset(size *types.Pos, cube_len int) (int, int) {
	x := (size.Z() * cube_len / 2) -
		(cube_len / 2)

	y := (size.X() * cube_len / 4) +
		(size.Y() * cube_len / 2) +
		(size.Z() * cube_len / 4) -
		cube_len

	return x, y
}

func GetIsoCubePosition(center_x, center_y, cube_len int, pos *types.Pos) (int, int) {
	x := center_x -
		(pos.Z() * cube_len / 2) +
		(pos.X() * cube_len / 2)

	y := center_y -
		(pos.X() * cube_len / 4) -
		(pos.Y() * cube_len / 2) -
		(pos.Z() * cube_len / 4)

	return x, y
}

func DrawIsoCube(img *image.RGBA, cube_len, x_offset, y_offset int, c1, c2, c3 color.Color) error {
	if cube_len%4 != 0 {
		return fmt.Errorf("cube_len must be divisible by 4")
	}
	if cube_len <= 4 {
		return fmt.Errorf("cube_len must be greater than 4")
	}

	half_len_zero_indexed := (cube_len / 2) - 1
	quarter_len := cube_len / 4

	// left/right part
	yo := 0
	for x := 0; x <= half_len_zero_indexed; x++ {
		for y := 0; y <= half_len_zero_indexed; y++ {
			// left
			img.Set(x_offset+x, y_offset+y+quarter_len+yo, c1)
			// right
			img.Set(x_offset+cube_len-1-x, y_offset+y+quarter_len+yo, c2)
		}
		if x%2 == 0 {
			yo = yo + 1
		}
	}

	// upper part
	yo = 0
	yl := 1
	for x := 0; x <= half_len_zero_indexed-1; x++ {
		for y := 0; y <= yl; y++ {
			// left
			img.Set(x_offset+1+x, y_offset+quarter_len-1-yo+y, c3)
			// right
			img.Set(x_offset+cube_len-2-x, y_offset+quarter_len-1-yo+y, c3)
		}
		if x%2 != 0 {
			yo = yo + 1
			yl = yl + 2
		}
	}

	return nil
}

func GetIsoNodeOrder(p *types.Pos) int {
	return (64000 - p.X()) + p.Y() + (64000 - p.Z())
}

func addAndClampUint8(a uint8, b int) uint8 {
	v := int(a) + b
	if v > 255 {
		return 255
	} else if v < 0 {
		return 0
	} else {
		return uint8(v)
	}
}

func ColorAdjust(c *color.RGBA, value int) *color.RGBA {
	return &color.RGBA{
		R: addAndClampUint8(c.R, value),
		G: addAndClampUint8(c.G, value),
		B: addAndClampUint8(c.B, value),
		A: c.A,
	}
}
