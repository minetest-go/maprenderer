package maprenderer

import (
	"fmt"

	"github.com/minetest-go/types"
)

func Probe(min, max, pos, ipos *types.Pos, na types.NodeAccessor, cr types.ColorResolver, skip_alpha bool) ([]*NodeWithColor, error) {
	nodes := []*NodeWithColor{}

	cpos := pos
	for cpos.IsWithin(min, max) {
		node, err := na(cpos)
		if err != nil {
			return nil, fmt.Errorf("getNode error @ %s: %v", cpos, err)
		}

		if node != nil && node.Name != "air" && node.Name != "ignore" {
			c := cr(node.Name, node.Param2)
			if c != nil {
				nodes = append(nodes, &NodeWithColor{
					Node:  node,
					Color: c,
				})

				if c.A == 255 || skip_alpha {
					// solid color or skip-param set
					break
				}
			}
		}
		cpos = cpos.Add(ipos)
	}

	return nodes, nil
}
