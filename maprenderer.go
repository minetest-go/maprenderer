package maprenderer

import (
	"image"

	"github.com/minetest-go/types"
)

type MapRenderOpts struct {
}

func NewDefaultMapRenderOpts() *MapRenderOpts {
	return &MapRenderOpts{}
}

func RenderMap(na types.NodeAccessor, cr types.ColorResolver, from, to *types.Pos, opts *MapRenderOpts) (*image.NRGBA, error) {
	// from = lowest, to = highest
	from, to = types.SortPos(from, to)
	search_dir := &types.Pos{0, -1, 0}

	// prepare image
	upLeft := image.Point{0, 0}
	lowRight := image.Point{to[0] - from[0] + 1, to[2] - from[2] + 1}
	img := image.NewNRGBA(image.Rectangle{upLeft, lowRight})

	for x := from.X(); x <= to.X(); x++ {
		for z := from.Z(); z <= to.Z(); z++ {
			// top-down search
			nodes, err := Probe(from, to, types.NewPos(x, to.Y(), z), search_dir, na, cr, false)
			if err != nil {
				return nil, err
			}

			if len(nodes) < 1 {
				continue
			}

			top_node := nodes[0]
			c := top_node.Color

			if len(nodes) == 1 && top_node.Color.A == 255 {
				// a single opaque node, draw with shadows

				// add shadows for view-blocking neighbors
				for _, above_pos := range []*types.Pos{{-1, 1, 0}, {0, 1, 1}} {
					nn, err := na(top_node.Pos.Add(above_pos))
					if err != nil {
						return nil, err
					}
					if nn != nil && cr(nn.Name, 0) != nil {
						c = ColorAdjust(c, -10)
					}
				}

				// lighten up if no nodes directly nearby
				for _, near_pos := range []*types.Pos{{-1, 0, 0}, {0, 0, 1}} {
					nn, err := na(top_node.Pos.Add(near_pos))
					if err != nil {
						return nil, err
					}
					if nn == nil || cr(nn.Name, 0) == nil {
						c = ColorAdjust(c, 10)
					}
				}
			} else {
				// multiple nodes with alpha channel
				c = nodes[len(nodes)-1].Color

				// bottom up color blending
				for i := len(nodes) - 2; i >= 0; i-- {
					node := nodes[i]
					c = BlendColor(c, node.Color, 3)
				}
			}

			img.Set(x-from.X(), to.Z()-z, *c)
		}
	}

	return img, nil
}
