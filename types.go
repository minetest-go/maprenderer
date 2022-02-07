package maprenderer

type Mapblock interface {
	GetNodeName(x, y, z int) string
	GetParam2(x, y, z int) int
	IsEmpty() bool
}

type MapblockPosGetter interface {
	GetX() int
	GetY() int
	GetZ() int
}

type MapblockPos struct {
	X, Y, Z int
}

func (pos *MapblockPos) GetX() int { return pos.X }
func (pos *MapblockPos) GetY() int { return pos.Y }
func (pos *MapblockPos) GetZ() int { return pos.Z }

type MapblockAccessor func(pos MapblockPosGetter) (Mapblock, error)
