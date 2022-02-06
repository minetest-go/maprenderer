package maprenderer

import (
	"errors"
	"image"

	"github.com/minetest-go/maprenderer/colormapping"
)

const (
	IMG_SCALE                         = 16
	IMG_SIZE                          = IMG_SCALE * 16
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

func (r *MapRenderer) Render(pos1, pos2 MapblockPos) (*image.NRGBA, error) {
	if pos1.X != pos2.X {
		return nil, errors.New("x does not line up")
	}

	if pos1.Z != pos2.Z {
		return nil, errors.New("z does not line up")
	}

	upLeft := image.Point{0, 0}
	lowRight := image.Point{IMG_SIZE, IMG_SIZE}
	img := image.NewNRGBA(image.Rectangle{upLeft, lowRight})

	maxY := pos1.Y
	minY := pos2.Y

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
			X: pos1.X,
			Y: mapBlockY,
			Z: pos1.Z,
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

					imgX := x * IMG_SCALE
					imgY := (15 - z) * IMG_SCALE

					r32, g32, b32, a32 := c.RGBA()
					r8, g8, b8, a8 := uint8(r32), uint8(g32), uint8(b32), uint8(a32)
					for Y := imgY; Y < imgY+IMG_SCALE; Y++ {
						ix := (Y*IMG_SIZE + imgX) << 2
						for X := 0; X < IMG_SCALE; X++ {
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
