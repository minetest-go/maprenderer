package maprenderer

import "math"

var tan30 = math.Tan(30 * math.Pi / 180)
var sqrt3div2 = 2 / math.Sqrt(3)

func GetIsoCubeSize(cubesize float64) (float64, float64) {
	return cubesize * 2, cubesize * sqrt3div2 * 2
}

func GetIsometricImageSize(size *Pos, cubesize float64) (int, int) {
	// max size of z or x axis
	max_xz := size.X()
	if size.Z() > max_xz {
		max_xz = size.Z()
	}

	size_x := cubesize * float64(size.X()+size.Z())
	size_y := cubesize * sqrt3div2 * float64(size.Y()+max_xz)

	return int(size_x), int(size_y)
}
