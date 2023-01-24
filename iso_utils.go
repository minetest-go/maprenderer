package maprenderer

import "math"

var tan30 = math.Tan(30 * math.Pi / 180)
var sqrt3div2 = 2 / math.Sqrt(3)
var sqrt2 = math.Sqrt(2)

func GetIsoCubeSize(cubesize float64) (float64, float64) {
	return cubesize * 2, cubesize * sqrt3div2 * 2
}

func GetIsometricImageSize(size *Pos, cubesize float64) (int, int) {
	// max size of z or x axis
	max_xz := size.X()
	if size.Z() > max_xz {
		max_xz = size.Z()
	}

	size_x := math.Ceil(cubesize * float64(size.X()+size.Z()+2))
	size_y := math.Ceil(cubesize * float64(size.Y()+max_xz+2) * sqrt3div2)

	return int(size_x), int(size_y)
}
