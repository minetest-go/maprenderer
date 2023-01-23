package maprenderer

import (
	"errors"
	"image"
	"image/color"
	"math"

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

func (r *IsoRenderer) Render(from, to *Pos) (image.Image, error) {
	// from = lowest, to = highest
	from, to = SortPos(from, to)
	direction := &Pos{1, -1, 1}

	// prepare image
	dc := gg.NewContext(600, 600) //TODO

	/*
		for y := from[1]; y <= to[1]; y++ {
			// right side
			for x := to[0]; x >= from[0]; x-- {
				err := r.renderPosition(dc, y-from[1], &Pos{x, y, from[2]}, from, direction)
				if err != nil {
					return nil, err
				}
			}

			// left side
			for z := to[2]; z >= from[2]; z-- {
				err := r.renderPosition(dc, y-from[1], &Pos{from[0], y, z}, from, direction)
				if err != nil {
					return nil, err
				}
			}
		}
	*/

	// top side
	for z := to.Z(); z >= from.Z(); z-- {
		for x := to.X(); x >= from.X(); x-- {
			iterations := max(to.X()-x, to.Z()-z)
			err := r.renderPosition(dc, iterations, &Pos{x, to.Y(), z}, from, direction)
			if err != nil {
				return nil, err
			}
		}
	}

	return dc.Image(), nil
}

func (r *IsoRenderer) renderPosition(dc *gg.Context, iterations int, pos, base_pos, direction *Pos) error {
	node, err := r.na.SearchNode(pos, direction, iterations)
	if err != nil {
		return err
	}

	if node == nil {
		// no node found or air
		return nil
	}

	c := r.cr(node.Name, node.Param2)
	if c != nil {
		rel_pos := node.Pos.Subtract(base_pos)
		r.drawBlock(dc, rel_pos, c)
	}

	return nil
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
