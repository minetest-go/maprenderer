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

type IsoDirection uint8

// look direction
const (
	// default direction, x and z increase
	IsoDirectionNorthEast IsoDirection = 0
	IsoDirectionNorthWest IsoDirection = 1
	IsoDirectionSouthWest IsoDirection = 2
	IsoDirectionSouthEast IsoDirection = 3
)

func NewIsoRenderer(cr ColorResolver, na NodeAccessor, height int) (*IsoRenderer, error) {
	if height%16 != 0 {
		return nil, errors.New("size is not a multiple of 16")
	}

	return &IsoRenderer{
		cr:     cr,
		na:     na,
		height: height,
		size:   10,
	}, nil
}

type IsoRenderer struct {
	cr     ColorResolver
	na     NodeAccessor
	height int
	size   float64
}

func (r *IsoRenderer) Render(from, to [3]int) (image.Image, error) {
	// from = lowest, to = highest
	from, to = SortPos(from, to)
	direction := [3]int{-1, -1, -1}

	// prepare image
	dc := gg.NewContext(600, 600) //TODO

	for y := from[1]; y <= to[1]; y++ {
		// right side
		for x := to[0]; x >= from[0]; x-- {
			err := r.renderPosition(dc, to[1]-y, [3]int{x, y, to[2]}, direction)
			if err != nil {
				return nil, err
			}
		}

		// left side
		for z := to[2]; z >= from[2]; z-- {
			err := r.renderPosition(dc, to[1]-y, [3]int{to[0], y, z}, direction)
			if err != nil {
				return nil, err
			}
		}
	}

	// top side
	for z := to[2]; z >= from[2]; z-- {
		for x := to[0]; x >= from[0]; x-- {
			err := r.renderPosition(dc, to[1]-from[1], [3]int{x, to[1], z}, direction)
			if err != nil {
				return nil, err
			}
		}
	}

	return dc.Image(), nil
}

func (r *IsoRenderer) renderPosition(dc *gg.Context, iterations int, pos, direction [3]int) error {
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
		r.drawBlock(dc, node.Pos, c)
	}

	return nil
}

func (r *IsoRenderer) getImagePos(x, y, z float64) (float64, float64) {
	xpos := (r.size * x) - (r.size * z)
	ypos := (r.size * tan30 * x) - (r.size * tan30 * z) - (r.size * sqrt3div2 * y)

	return xpos + 50, ypos + 200
}

func (r *IsoRenderer) drawBlock(dc *gg.Context, pos [3]int, color *color.RGBA) {
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
