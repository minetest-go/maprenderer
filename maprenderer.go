package maprenderer

import (
	"image"
)

func RenderMap(from, to *Pos, na NodeAccessor, cr ColorResolver) (*image.NRGBA, error) {
	// from = lowest, to = highest
	from, to = SortPos(from, to)
	search_dir := &Pos{0, -1, 0}

	// prepare image
	upLeft := image.Point{0, 0}
	lowRight := image.Point{to[0] - from[0] + 1, to[2] - from[2] + 1}
	img := image.NewNRGBA(image.Rectangle{upLeft, lowRight})

	for x := from.X(); x <= to.X(); x++ {
		for z := from.Z(); z <= to.Z(); z++ {
			// top-down search
			node, err := na.SearchNode(&Pos{x, to.Y(), z}, search_dir, [2]*Pos{from, to})
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

			// add shadows for view-blocking neighbors
			for _, above_pos := range []*Pos{{-1, 1, 0}, {0, 1, 1}} {
				nn, err := na.GetNode(node.Pos.Add(above_pos))
				if err != nil {
					return nil, err
				}
				if nn != nil && cr(nn.Name, 0) != nil {
					c = AddColorComponent(c, -10)
				}
			}

			// lighten up if no nodes directly nearby
			for _, near_pos := range []*Pos{{-1, 0, 0}, {0, 0, 1}} {
				nn, err := na.GetNode(node.Pos.Add(near_pos))
				if err != nil {
					return nil, err
				}
				if nn == nil || cr(nn.Name, 0) == nil {
					c = AddColorComponent(c, 10)
				}
			}

			img.Set(x, to.Z()-z, *c)
		}
	}

	return img, nil
}
