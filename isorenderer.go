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
		CubeLen:            8,
		EnableTransparency: false,
	}
}

func RenderIsometric(na types.NodeAccessor, cr types.ColorResolver, from, to *types.Pos, opts *IsoRenderOpts) (image.Image, error) {
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

	// right layer
	for x := min.X(); x <= max.X(); x++ {
		for y := min.Y(); y <= max.Y()-1; y++ {
			pnodes, err := Probe(min, max, types.NewPos(x, y, min.Z()), ipos, na, cr, skip_alpha)
			if err != nil {
				return nil, fmt.Errorf("probe error, right layer: %v", err)
			}
			nodes = append(nodes, pnodes...)
		}
	}

	// left layer
	for z := min.Z() + 1; z <= max.Z(); z++ {
		for y := min.Y(); y <= max.Y()-1; y++ {
			pnodes, err := Probe(min, max, types.NewPos(min.X(), y, z), ipos, na, cr, skip_alpha)
			if err != nil {
				return nil, fmt.Errorf("probe error, left layer: %v", err)
			}
			nodes = append(nodes, pnodes...)
		}
	}

	slices.SortFunc(nodes, SortNodesWithColor)

	for _, n := range nodes {
		rel_pos := n.Pos.Subtract(min)
		c1 := ColorAdjust(n.Color, 0)
		if !opts.EnableTransparency {
			// disable alpha channel
			c1.A = 255
		}

		// brighter/darker colors for the other sides
		c2 := ColorAdjust(c1, -10)
		c3 := ColorAdjust(c1, 10)

		x, y := GetIsoCubePosition(center_x, center_y, opts.CubeLen, rel_pos)
		err := DrawIsoCube(img, opts.CubeLen, x, y, c1, c2, c3)
		if err != nil {
			return nil, fmt.Errorf("DrawIsoCube error: %v", err)
		}
	}

	return img, nil
}
