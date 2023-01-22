package maprenderer

import "image/color"

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

// ng stuff

// simple node-definition
type Node struct {
	Pos    [3]int
	Name   string `json:"name"`
	Param1 int    `json:"param1"`
	Param2 int    `json:"param2"`
}

type NodeAccessor interface {
	// returns the first non-air, non-ignore node in the search direction, nil if none found
	SearchNode(pos, direction [3]int, iterations int) (*Node, error)
	// returns the node at the given position, nil if no node found
	GetNode(pos [3]int) (*Node, error)
}

// resolves the node-name and param2 to a color, nil if no color-mapping found
type ColorResolver func(name string, param2 int) *color.RGBA
