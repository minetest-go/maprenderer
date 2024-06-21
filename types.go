package maprenderer

import "image/color"

// simple node-definition
type Node struct {
	Pos    *Pos
	Name   string `json:"name"`
	Param1 int    `json:"param1"`
	Param2 int    `json:"param2"`
}

type NodeWithColor struct {
	*Node
	Color *color.RGBA
}

// returns the node at the given position, nil if no node found
type NodeAccessor func(pos *Pos) (*Node, error)

// resolves the node-name and param2 to a color, nil if no color-mapping found
type ColorResolver func(name string, param2 int) *color.RGBA
