package maprenderer

import (
	"errors"
	"image"

	"github.com/minetest-go/maprenderer/colormapping"
)

const (
	EXPECTED_BLOCKS_PER_FLAT_MAPBLOCK = 16 * 16
)

func NewMapRenderer(cm *colormapping.ColorMapping, mba MapblockAccessor, size int) (*MapRenderer, error) {
	if size%16 != 0 {
		return nil, errors.New("size is not a multiple of 16")
	}

	return &MapRenderer{
		cm:    cm,
		mba:   mba,
		size:  size,
		scale: int(size / 16),
	}, nil
}

type MapRenderer struct {
	cm    *colormapping.ColorMapping
	mba   MapblockAccessor
	size  int
	scale int
}

func (r *MapRenderer) IsViewBocking(nodename string) bool {
	if nodename == "air" || nodename == "" {
		// not visible
		return false
	}

	c := r.cm.GetColor(nodename, 0)
	return c != nil // has color == view-blocking
}

func (r *MapRenderer) Render(pos MapblockPosGetter, y_block_height int) (*image.NRGBA, error) {

	upLeft := image.Point{0, 0}
	lowRight := image.Point{r.size, r.size}
	img := image.NewNRGBA(image.Rectangle{upLeft, lowRight})

	maxY := pos.GetY() + y_block_height
	minY := pos.GetY()

	if minY > maxY {
		maxY, minY = minY, maxY
	}

	foundBlocks := 0
	xzOccupationMap := make([][]bool, 16)
	for x := range xzOccupationMap {
		xzOccupationMap[x] = make([]bool, 16)
	}

	for mapBlockY := maxY; mapBlockY >= minY; mapBlockY-- {
		currentPos := MapblockPos{
			X: pos.GetX(),
			Y: mapBlockY,
			Z: pos.GetZ(),
		}

		mb, err := r.mba(&currentPos)

		if err != nil {
			return nil, err
		}

		if mb == nil || mb.IsEmpty() {
			continue
		}

		for x := 0; x < 16; x++ {
			for z := 0; z < 16; z++ {
				for y := 15; y >= 0; y-- {
					if xzOccupationMap[x][z] {
						break
					}

					nodeName := mb.GetNodeName(x, y, z)
					param2 := mb.GetParam2(x, y, z)

					if nodeName == "" {
						continue
					}

					c := r.cm.GetColor(nodeName, param2)

					if c == nil {
						continue
					}

					var left, leftAbove, top, topAbove string

					if x > 0 {
						//same mapblock
						left = mb.GetNodeName(x-1, y, z)
						if y < 15 {
							leftAbove = mb.GetNodeName(x-1, y+1, z)
						}

					} else {
						//neighbouring mapblock
						neighbourPos := MapblockPos{
							X: currentPos.X - 1,
							Y: currentPos.Y,
							Z: currentPos.Z,
						}
						neighbourMapblock, err := r.mba(&neighbourPos)

						if neighbourMapblock != nil && err == nil {
							left = neighbourMapblock.GetNodeName(15, y, z)
							if y < 15 {
								leftAbove = neighbourMapblock.GetNodeName(15, y+1, z)
							}
						}
					}

					if z < 14 {
						//same mapblock
						top = mb.GetNodeName(x, y, z+1)
						if y < 15 {
							topAbove = mb.GetNodeName(x, y+1, z+1)
						}

					} else {
						//neighbouring mapblock
						neighbourPos := MapblockPos{
							X: currentPos.X,
							Y: currentPos.Y,
							Z: currentPos.Z + 1,
						}
						neighbourMapblock, err := r.mba(&neighbourPos)

						if neighbourMapblock != nil && err == nil {
							top = neighbourMapblock.GetNodeName(x, y, 0)
							if y < 15 {
								topAbove = neighbourMapblock.GetNodeName(x, y+1, 0)
							}
						}
					}

					if r.IsViewBocking(leftAbove) {
						//add shadow
						c = AddColorComponent(c, -10)
					}

					if r.IsViewBocking(topAbove) {
						//add shadow
						c = AddColorComponent(c, -10)
					}

					if !r.IsViewBocking(left) {
						//add light
						c = AddColorComponent(c, 10)
					}

					if !r.IsViewBocking(top) {
						//add light
						c = AddColorComponent(c, 10)
					}

					imgX := x * r.scale
					imgY := (15 - z) * r.scale

					r32, g32, b32, a32 := c.RGBA()
					r8, g8, b8, a8 := uint8(r32), uint8(g32), uint8(b32), uint8(a32)
					for Y := imgY; Y < imgY+r.scale; Y++ {
						ix := (Y*r.size + imgX) << 2
						for X := 0; X < r.scale; X++ {
							img.Pix[ix] = r8
							ix++
							img.Pix[ix] = g8
							ix++
							img.Pix[ix] = b8
							ix++
							img.Pix[ix] = a8
							ix++
						}
					}

					//not transparent, mark as rendered
					foundBlocks++
					xzOccupationMap[x][z] = true

					if foundBlocks == EXPECTED_BLOCKS_PER_FLAT_MAPBLOCK {
						return img, nil
					}
				}
			}
		}
	}

	return nil, nil
}
