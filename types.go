package maprenderer

import "image/color"

// simple node-definition
type Node struct {
	Pos    *Pos
	Name   string `json:"name"`
	Param1 int    `json:"param1"`
	Param2 int    `json:"param2"`
}

type NodeAccessor interface {
	// returns the first non-air, non-ignore node in the search direction, nil if none found
	// start and end are inclusive
	SearchNode(start, direction *Pos, bounds [2]*Pos) (*Node, error)
	// returns the node at the given position, nil if no node found
	GetNode(pos *Pos) (*Node, error)
}

// resolves the node-name and param2 to a color, nil if no color-mapping found
type ColorResolver func(name string, param2 int) *color.RGBA
