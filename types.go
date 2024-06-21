package maprenderer

import (
	"image/color"

	"github.com/minetest-go/types"
)

type NodeWithColor struct {
	*types.Node
	Color *color.RGBA
}

// returns the node at the given position, nil if no node found
type NodeAccessor func(pos *types.Pos) (*types.Node, error)

// resolves the node-name and param2 to a color, nil if no color-mapping found
type ColorResolver func(name string, param2 int) *color.RGBA
