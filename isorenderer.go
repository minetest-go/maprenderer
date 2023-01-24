package maprenderer

import (
	"errors"
	"image"
	"image/color"
	"math"
	"sort"

	"github.com/fogleman/gg"
)

var tan30 = math.Tan(30 * math.Pi / 180)
var sqrt3div2 = 2 / math.Sqrt(3)

func NewIsoRenderer(cr ColorResolver, na NodeAccessor, height int) (*IsoRenderer, error) {
	if height%16 != 0 {
		return nil, errors.New("size is not a multiple of 16")
	}

	return &IsoRenderer{
		cr:     cr,
		na:     na,
		height: height,
		size:   6,
	}, nil
}

type IsoRenderer struct {
	cr     ColorResolver
	na     NodeAccessor
	height int
	size   float64
}

type IsometricNode struct {
	*Node
	*color.RGBA
	Order int
}

func (r *IsoRenderer) Render(from, to *Pos) (image.Image, error) {
	// from = lowest, to = highest
	from, to = SortPos(from, to)
	direction := &Pos{1, -1, 1}

	nodes := make([]*IsometricNode, 0)

	// prepare image
	dc := gg.NewContext(600, 600) //TODO

	for y := from.Y(); y <= to.Y(); y++ {
		// right side
		for x := to[0]; x >= from[0]; x-- {
			n, err := r.searchNode(&Pos{x, y, from[2]}, direction, from, [2]*Pos{from, to})
			if err != nil {
				return nil, err
			}
			if n != nil {
				nodes = append(nodes, n)
			}
		}

		// left side
		for z := to[2]; z >= from[2]; z-- {
			n, err := r.searchNode(&Pos{from[0], y, z}, direction, from, [2]*Pos{from, to})
			if err != nil {
				return nil, err
			}
			if n != nil {
				nodes = append(nodes, n)
			}
		}
	}

	// top side
	for z := to.Z(); z >= from.Z(); z-- {
		for x := to.X(); x >= from.X(); x-- {
			n, err := r.searchNode(&Pos{x, to.Y(), z}, direction, from, [2]*Pos{from, to})
			if err != nil {
				return nil, err
			}
			if n != nil {
				nodes = append(nodes, n)
			}
		}
	}

	sort.Slice(nodes, func(i int, j int) bool {
		return nodes[i].Order < nodes[j].Order
	})

	for _, node := range nodes {
		r.drawBlock(dc, node.Pos, node.RGBA)
	}

	return dc.Image(), nil
}

func (r *IsoRenderer) searchNode(pos, direction, base_pos *Pos, bounds [2]*Pos) (*IsometricNode, error) {
	node, err := r.na.SearchNode(pos, direction, bounds)
	if err != nil {
		return nil, err
	}

	if node == nil {
		// no node found or air
		return nil, nil
	}

	c := r.cr(node.Name, node.Param2)
	if c != nil {
		return &IsometricNode{
			Node:  node,
			RGBA:  c,
			Order: node.Pos.Y() + ((bounds[1].X() - node.Pos.X()) * bounds[1].X()) + ((bounds[1].Z() - node.Pos.Z()) + bounds[1].Z()),
		}, nil
	}

	return nil, nil
}

func (r *IsoRenderer) getImagePos(x, y, z float64) (float64, float64) {
	xpos := 300 + (r.size * x) - (r.size * z)
	ypos := 450 - (r.size * tan30 * x) - (r.size * tan30 * z) - (r.size * sqrt3div2 * y)

	return xpos, ypos
}

func (r *IsoRenderer) drawBlock(dc *gg.Context, pos *Pos, color *color.RGBA) {
	x, y := r.getImagePos(float64(pos[0]), float64(pos[1]), float64(pos[2]))
	radius := r.size

	// right side
	dc.MoveTo(radius+x, (radius*tan30)+y)
	dc.LineTo(x, (radius*sqrt3div2)+y)
	dc.LineTo(x, y)
	dc.LineTo(radius+x, -(radius*tan30)+y)
	dc.ClosePath()
	dc.SetRGB255(int(color.R), int(color.G), int(color.B))
	dc.Fill()

	// left side
	dc.MoveTo(x, (radius*sqrt3div2)+y)
	dc.LineTo(-radius+x, (radius*tan30)+y)
	dc.LineTo(-radius+x, -(radius*tan30)+y)
	dc.LineTo(x, y)
	dc.ClosePath()
	AdjustAndFill(dc, int(color.R), int(color.G), int(color.B), -20)
	dc.Fill()

	// top side
	dc.MoveTo(-radius+x, -(radius*tan30)+y)
	dc.LineTo(x, -(radius*sqrt3div2)+y)
	dc.LineTo(radius+x, -(radius*tan30)+y)
	dc.LineTo(x, y)
	dc.ClosePath()
	AdjustAndFill(dc, int(color.R), int(color.G), int(color.B), 20)
	dc.Fill()
}
