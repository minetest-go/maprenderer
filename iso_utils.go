package maprenderer

func GetIsoCubeSize(size float64) (float64, float64) {
	return size * 2, size * sqrt3div2 * 2
}

func GetIsometricImageSize(from, to *Pos, cubesize int) (int, int) {
	cx_s, cy_s := GetIsoCubeSize(float64(cubesize))

	// max size of z or x axis
	x_diff := to.X() - from.X()
	y_diff := to.Y() - from.Y()
	z_diff := to.Z() - from.Z()
	max_xz := x_diff
	if z_diff > x_diff {
		max_xz = z_diff
	}

	size_x := int(cx_s) * (x_diff + z_diff) / 2
	size_y := int(cy_s) * (y_diff + max_xz) / 2

	return size_x, size_y
}
