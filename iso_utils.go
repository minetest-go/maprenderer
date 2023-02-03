package maprenderer

import "math"

var tan30 = math.Tan(30 * math.Pi / 180)
var sqrt3div2 = 2 / math.Sqrt(3)
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

	x_pos := (float64(rel_pos.X()) * cube_x / 2) +
		(float64(rel_pos.Z()) * cube_x / 2)

	y_pos := (float64(size.Y()-rel_pos.Y()-1) * cube_y / 2) +
		(float64(size.X()-rel_pos.X()-1) * cube_y / 2) +
		(float64(size.Z()-rel_pos.Z()-1) * cube_y / 2)

	return x_pos, y_pos
}
