package maprenderer

import (
	"image"
)

type MapRenderOpts struct {
}

func RenderMap(na NodeAccessor, cr ColorResolver, from, to *Pos, opts *MapRenderOpts) (*image.NRGBA, error) {
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
			nodes, err := Probe(from, to, NewPos(x, to.Y(), z), search_dir, na, cr, true)
			if err != nil {
				return nil, err
			}

			if len(nodes) < 1 {
				continue
			}

			node := nodes[0]

			c := cr(node.Name, node.Param2)
			if c == nil {
				continue
			}

			// add shadows for view-blocking neighbors
			for _, above_pos := range []*Pos{{-1, 1, 0}, {0, 1, 1}} {
				nn, err := na(node.Pos.Add(above_pos))
				if err != nil {
					return nil, err
				}
				if nn != nil && cr(nn.Name, 0) != nil {
					c = ColorAdjust(c, -10)
				}
			}

			// lighten up if no nodes directly nearby
			for _, near_pos := range []*Pos{{-1, 0, 0}, {0, 0, 1}} {
				nn, err := na(node.Pos.Add(near_pos))
				if err != nil {
					return nil, err
				}
				if nn == nil || cr(nn.Name, 0) == nil {
					c = ColorAdjust(c, 10)
				}
			}

			img.Set(x-from.X(), to.Z()-z, *c)
		}
	}

	return img, nil
}
