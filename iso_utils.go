package maprenderer

import "math"

var tan30 = math.Tan(30 * math.Pi / 180)
var sqrt3div2 = 2 / math.Sqrt(3)

// returns the outer cube dimensions (of the surrounding rectangle)
// "cubesize" is the length from the center to the right or left side (2*cubesize is the width)
func GetIsoCubeSize(cubesize float64) (float64, float64) {
	return cubesize * 2, cubesize * sqrt3div2 * 2
}

func GetIsometricImageSize(size *Pos, cubesize float64) (int, int) {
	cube_x, cube_y := GetIsoCubeSize(cubesize)

	// max size of z or x axis
	max_xz := size.X()
	if size.Z() > max_xz {
		max_xz = size.Z()
	}

	size_x := math.Floor(cube_x * float64(size.X()+size.Z()) / 2)
	size_y := math.Floor(cube_y * float64(size.Y()+max_xz) / 2)

	return int(size_x), int(size_y)
}

// returns the left/top aligned image position for the cube at given position
func GetImagePos(rel_pos *Pos, size_x, size_y int, cubesize float64) (float64, float64) {
	// floating point coords
	x := float64(rel_pos.X())
	y := float64(rel_pos.Y())
	z := float64(rel_pos.Z())

	// max size of z or x axis
	max_xz := x
	if z > max_xz {
		max_xz = z
	}

	xpos := ((cubesize * x) - (cubesize * z)) + (float64(size_x) / 2) - cubesize
	ypos := float64(size_y) - (cubesize * sqrt3div2 * 2) -
		(cubesize * tan30 * x) -
		(cubesize * tan30 * z) -
		(cubesize * sqrt3div2 * y)

	return xpos, ypos
}
