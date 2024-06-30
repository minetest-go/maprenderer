package maprenderer

import (
	"fmt"
	"image"
	"slices"

	"github.com/minetest-go/types"
)

type IsoTileRenderOpts struct {
	CubeLen int
}

func NewDefaultIsoTileRenderOpts() *IsoTileRenderOpts {
	return &IsoTileRenderOpts{
		CubeLen: 8,
	}
}

func RenderIsometricTile(na types.NodeAccessor, cr types.ColorResolver, from, to *types.Pos, opts *IsoTileRenderOpts) (image.Image, error) {
	if opts == nil {
		opts = NewDefaultIsoTileRenderOpts()
	}

	min, max := types.SortPos(from, to)
	size := to.Subtract(from).Add(types.NewPos(1, 1, 1))
	top_size := types.NewPos(size.X(), 1, size.Z())
	probe_bounds_min := types.NewPos(min.X()-(size.Y()*2), min.Y(), min.Z()-(size.Y()*2))
	probe_bounds_max := types.NewPos(max.X()+(size.Y()*2), max.Y(), max.Z()+(size.Y()*2))

	width, height := GetIsometricImageSize(top_size, opts.CubeLen)
	center_x, center_y := GetIsoCenterCubeOffset(size, opts.CubeLen)
	img := image.NewRGBA(image.Rectangle{Min: image.Point{}, Max: image.Point{X: width, Y: height}})

	ipos := types.NewPos(1, -1, 1)

	nodes := []*NodeWithColor{}

	// top layer
	for x := min.X(); x <= max.X(); x++ {
		for z := min.Z(); z <= max.Z(); z++ {
			pnodes, err := Probe(min, max, types.NewPos(x, max.Y(), z), ipos, na, cr, true)
			if err != nil {
				return nil, fmt.Errorf("probe error, top layer: %v", err)
			}
			nodes = append(nodes, pnodes...)
		}
	}

	// bottom right corner
	for zi := 0; zi <= size.Z(); zi++ {
		for x := min.X() + zi - 3; x <= max.X()-zi+1; x++ {
			pos := types.NewPos(x, max.Y(), min.Z()-zi)
			pnodes, err := Probe(probe_bounds_min, probe_bounds_max, pos, ipos, na, cr, true)
			if err != nil {
				return nil, fmt.Errorf("probe error, bottom right corner: %v", err)
			}
			nodes = append(nodes, pnodes...)
		}
	}

	// bottom left corner
	for xi := 0; xi <= size.X(); xi++ {
		for z := min.Z() + xi - 3; z <= max.Z()-xi+1; z++ {
			pos := types.NewPos(min.X()-xi, max.Y(), z)
			pnodes, err := Probe(probe_bounds_min, probe_bounds_max, pos, ipos, na, cr, true)
			if err != nil {
				return nil, fmt.Errorf("probe error, bottom left corner: %v", err)
			}
			nodes = append(nodes, pnodes...)
		}
	}

	slices.SortFunc(nodes, SortNodesWithColor)

	for _, n := range nodes {
		rel_pos := n.Pos.Subtract(min)
		c1 := ColorAdjust(n.Color, 0)
		// disable alpha channel
		c1.A = 255

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
