package maprenderer

import (
	"image/color"
	"math"

	"github.com/fogleman/gg"
)

var sqrt3 = math.Sqrt(3)
var sin30 = math.Sin(30 * math.Pi / 180)

// returns the outer cube dimensions (of the surrounding rectangle)
// "cubesize" is the true sidelength (https://en.wikipedia.org/wiki/Isometric_projection#/media/File:3D_shapes_in_isometric_projection.svg)
func GetIsoCubeSize(cubesize float64) (float64, float64) {
	return cubesize * sqrt3, cubesize * 2
}

func GetIsometricImageSize(size *Pos, cubesize float64) (int, int) {
	cube_x, cube_y := GetIsoCubeSize(cubesize)

	// max size of z or x axis
	max_xz := size.X()
	if size.Z() > max_xz {
		max_xz = size.Z()
	}

	size_x := math.Ceil(cube_x * float64(size.X()+size.Z()) / 2)
	size_y := math.Ceil(cube_y * float64(size.Y()+max_xz) / 2)

	return int(size_x), int(size_y)
}

// returns the left/top aligned image position for the cube at given position
func GetImagePos(rel_pos, size *Pos, size_x, size_y int, cubesize float64) (float64, float64) {
	// floating point coords
	cube_x, cube_y := GetIsoCubeSize(cubesize)

	x_pos := (float64(rel_pos.X()) * cube_x / 2) -
		(float64(rel_pos.Z()) * cube_x / 2) + 140

	y_pos := float64(size_y) -
		(float64(rel_pos.Y()) * cube_y / 2) -
		(float64(rel_pos.X()) * cube_y / 4) -
		(float64(rel_pos.Z()) * cube_y / 4)

	return x_pos, y_pos
}

func DrawCube(dc *gg.Context, c *color.RGBA, size float64, offset_x, offset_y float64) {
	size_x, size_y := GetIsoCubeSize(size)

	// center position
	center_x := (size_x / 2) + offset_x
	center_y := (size_y / 2) + offset_y

	// calculate ends
	end_x := offset_x + size_x
	end_y := offset_y + size_y

	// proportional size
	sin30_proportional := sin30 * size

	// right side
	dc.SetRGBA255(int(c.R), int(c.G), int(c.B), int(c.A))
	dc.MoveTo(center_x, center_y)
	dc.LineTo(end_x, center_y-sin30_proportional)
	dc.LineTo(end_x, end_y-sin30_proportional)
	dc.LineTo(center_x, end_y)
	dc.ClosePath()
	dc.Fill()

	// left side
	dc.SetRGBA255(
		AdjustColorComponent(c.R, -20),
		AdjustColorComponent(c.G, -20),
		AdjustColorComponent(c.B, -20),
		int(c.A),
	)
	dc.MoveTo(center_x, center_y)
	dc.LineTo(center_x, end_y)
	dc.LineTo(offset_x, end_y-sin30_proportional)
	dc.LineTo(offset_x, center_y-sin30_proportional)
	dc.ClosePath()
	dc.Fill()

	// top side
	dc.SetRGBA255(
		AdjustColorComponent(c.R, 20),
		AdjustColorComponent(c.G, 20),
		AdjustColorComponent(c.B, 20),
		int(c.A),
	)
	dc.MoveTo(center_x, center_y)
	dc.LineTo(offset_x, center_y-sin30_proportional)
	dc.LineTo(center_x, offset_y)
	dc.LineTo(end_x, center_y-sin30_proportional)
	dc.ClosePath()
	dc.Fill()
}
