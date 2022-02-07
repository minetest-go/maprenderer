package maprenderer

import (
	"errors"
	"image"
	"image/color"
	"math"

	"github.com/fogleman/gg"
	"github.com/minetest-go/maprenderer/colormapping"
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

func NewIsoRenderer(cm *colormapping.ColorMapping, mba MapblockAccessor, height int) (*IsoRenderer, error) {
	if height%16 != 0 {
		return nil, errors.New("size is not a multiple of 16")
	}

	return &IsoRenderer{
		cm:     cm,
		mba:    mba,
		height: height,
		scale:  int(height / 16),
		size:   10, //TODO
	}, nil
}

type IsoRenderer struct {
	cm     *colormapping.ColorMapping
	mba    MapblockAccessor
	height int
	scale  int
	size   float64
}

func (r *IsoRenderer) Render(pos MapblockPosGetter, y_block_height int, direction IsoDirection) (*image.NRGBA, error) {
	// stub
	return nil, nil
}

func (r *IsoRenderer) GetImagePos(x, y, z float64) (float64, float64) {
	xpos := (r.size * x) - (r.size * z)
	ypos := (r.size * tan30 * x) - (r.size * tan30 * z) - (r.size * sqrt3div2 * y)

	return xpos, ypos
}

func (r *IsoRenderer) drawBlock(dc *gg.Context, color *color.RGBA) {
	//x, y := r.GetImagePos(float64(block.X), float64(block.Y), float64(block.Z))
	x, y := 0.0, 0.0
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
