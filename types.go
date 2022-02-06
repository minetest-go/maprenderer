package maprenderer

type Mapblock interface {
	GetNodeName(x, y, z int) string
	GetParam2(x, y, z int) int
	IsEmpty() bool
}

type MapblockPos struct {
	X int
	Y int
	Z int
}

type MapblockAccessor func(pos *MapblockPos) (Mapblock, error)
