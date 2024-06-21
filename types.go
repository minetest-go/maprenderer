package maprenderer

import (
	"image/color"

	"github.com/minetest-go/types"
)

type NodeWithColor struct {
	*types.Node
	Color *color.RGBA
}
