package maprenderer

import (
	"fmt"
	"image"
	"slices"

	"github.com/minetest-go/types"
)

type IsoRenderOpts struct {
	CubeLen            int
	EnableTransparency bool
}

func NewDefaultIsoRenderOpts() *IsoRenderOpts {
	return &IsoRenderOpts{
		CubeLen:            4,
		EnableTransparency: false,
	}
}

func RenderIsometric(na NodeAccessor, cr ColorResolver, from, to *types.Pos, opts *IsoRenderOpts) (image.Image, error) {
	if opts == nil {
		opts = NewDefaultIsoRenderOpts()
	}

	min, max := types.SortPos(from, to)
	size := to.Subtract(from).Add(types.NewPos(1, 1, 1))

	width, height := GetIsometricImageSize(size, opts.CubeLen)
	center_x, center_y := GetIsoCenterCubeOffset(size, opts.CubeLen)
	img := image.NewRGBA(image.Rectangle{Min: image.Point{}, Max: image.Point{X: width, Y: height}})

	skip_alpha := !opts.EnableTransparency
	ipos := types.NewPos(1, -1, 1)

	nodes := []*NodeWithColor{}

	// top layer
	for x := min.X(); x <= max.X(); x++ {
		for z := min.Z(); z <= max.Z(); z++ {
			pnodes, err := Probe(min, max, types.NewPos(x, max.Y(), z), ipos, na, cr, skip_alpha)
			if err != nil {
				return nil, fmt.Errorf("probe error, top layer: %v", err)
			}
			nodes = append(nodes, pnodes...)
		}
	}

	// left layer (without top stride)
	for x := min.X(); x <= max.X(); x++ {
		for y := min.Y(); y < max.Y()-1; y++ {
			pnodes, err := Probe(min, max, types.NewPos(x, y, min.Z()), ipos, na, cr, skip_alpha)
			if err != nil {
				return nil, fmt.Errorf("probe error, left layer: %v", err)
			}
			nodes = append(nodes, pnodes...)
		}
	}

	// right layer (without top and left stride)
	for z := min.Z() + 1; z < max.Z(); z++ {
		for y := min.Y(); y <= max.Y()-1; y++ {
			pnodes, err := Probe(min, max, types.NewPos(min.X(), y, z), ipos, na, cr, skip_alpha)
			if err != nil {
				return nil, fmt.Errorf("probe error, right layer: %v", err)
			}
			nodes = append(nodes, pnodes...)
		}
	}

	rel_max := max.Subtract(min)
	slices.SortFunc(nodes, func(n1, n2 *NodeWithColor) int {
		o1 := GetIsoNodeOrder(n1.Pos.Subtract(min), rel_max)
		o2 := GetIsoNodeOrder(n2.Pos.Subtract(min), rel_max)
		return o1 - o2
	})

	for _, n := range nodes {
		rel_pos := n.Pos.Subtract(min)
		c2 := ColorAdjust(n.Color, -10)
		c3 := ColorAdjust(n.Color, 10)

		x, y := GetIsoCubePosition(center_x, center_y, opts.CubeLen, rel_pos)
		err := DrawIsoCube(img, opts.CubeLen, x, y, n.Color, c2, c3)
		if err != nil {
			return nil, fmt.Errorf("DrawIsoCube error: %v", err)
		}
	}

	return img, nil
}
