package maprenderer

import (
	"image"
)

func RenderMap(from, to [3]int, na NodeAccessor, cr ColorResolver) (*image.NRGBA, error) {
	// from = lowest, to = highest
	from, to = SortPos(from, to)
	y_diff := to[1] - from[1]
	search_dir := [3]int{0, -1, 0}

	// prepare image
	upLeft := image.Point{0, 0}
	lowRight := image.Point{to[0] - from[0] + 1, to[2] - from[2] + 1}
	img := image.NewNRGBA(image.Rectangle{upLeft, lowRight})

	for x := from[0]; x <= to[0]; x++ {
		for z := from[2]; z <= to[2]; z++ {
			// top-down search
			node, err := na.SearchNode([3]int{x, to[1], z}, search_dir, y_diff)
			if err != nil {
				return nil, err
			}

			if node == nil {
				continue
			}

			c := cr(node.Name, node.Param2)
			if c == nil {
				continue
			}

			img.Set(x, img.Rect.Dx()-z, *c)
		}
	}

	return img, nil
}
