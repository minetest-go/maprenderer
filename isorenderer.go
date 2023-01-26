package maprenderer

import (
	"image"
	"image/color"
	"image/draw"
	"sort"
)

func NewIsoRenderer(cr ColorResolver, na NodeAccessor, cubesize int) (*IsoRenderer, error) {
	return &IsoRenderer{
		cr:       cr,
		na:       na,
		rc:       NewIsoRenderCache(),
		cubesize: float64(cubesize),
	}, nil
}

type IsoRenderer struct {
	cr       ColorResolver
	na       NodeAccessor
	rc       *IsoRenderCache
	cubesize float64
}

type IsometricNode struct {
	*Node
	*color.RGBA
	Order int
}

func (r *IsoRenderer) Render(from, to *Pos) (image.Image, error) {
	// from = lowest, to = highest
	from, to = SortPos(from, to)
	direction := &Pos{1, -1, 1}

	nodes := make([]*IsometricNode, 0)

	// skip top layer (drawn later)
	for y := from.Y(); y < to.Y(); y++ {
		// right side (skip an already drawn row)
		for x := to.X(); x > from.X(); x-- {
			n, err := r.searchNode(&Pos{x, y, from.Z()}, direction, from, [2]*Pos{from, to})
			if err != nil {
				return nil, err
			}
			if n != nil {
				nodes = append(nodes, n)
			}
		}

		// left side
		for z := to.Z(); z >= from.Z(); z-- {
			n, err := r.searchNode(&Pos{from.X(), y, z}, direction, from, [2]*Pos{from, to})
			if err != nil {
				return nil, err
			}
			if n != nil {
				nodes = append(nodes, n)
			}
		}
	}

	// top side
	for z := to.Z(); z >= from.Z(); z-- {
		for x := to.X(); x >= from.X(); x-- {
			n, err := r.searchNode(&Pos{x, to.Y(), z}, direction, from, [2]*Pos{from, to})
			if err != nil {
				return nil, err
			}
			if n != nil {
				nodes = append(nodes, n)
			}
		}
	}

	sort.Slice(nodes, func(i int, j int) bool {
		return nodes[i].Order < nodes[j].Order
	})

	// prepare image
	//dc := gg.NewContext(600, 600) //TODO

	size := to.Subtract(from).Add(NewPos(1, 1, 1))
	size_x, size_y := GetIsometricImageSize(size, r.cubesize)
	img := image.NewRGBA(image.Rect(0, 0, size_x, size_y))

	for _, node := range nodes {
		rel_pos := node.Pos.Subtract(from)
		x, y := r.getImagePos(float64(rel_pos.X()), float64(rel_pos.Y()), float64(rel_pos.Z()), size_x, size_y)

		cube_img := r.rc.GetCachedIsoCubeImage(node.RGBA, r.cubesize)
		p1 := image.Point{X: int(x), Y: int(y)}
		r := image.Rectangle{
			p1, p1.Add(cube_img.Bounds().Size()),
		}

		// NOTE: the native "draw.Draw" function doesn't work with transparency
		draw.Draw(img, r, cube_img, image.Point{0, 0}, draw.Over)
		//dc.DrawImage(cube_img, int(math.Floor(x)), int(math.Floor(y)))
	}

	return img, nil
}

func (r *IsoRenderer) searchNode(pos, direction, base_pos *Pos, bounds [2]*Pos) (*IsometricNode, error) {
	node, err := r.na.SearchNode(pos, direction, bounds)
	if err != nil {
		return nil, err
	}

	if node == nil {
		// no node found or air
		return nil, nil
	}

	c := r.cr(node.Name, node.Param2)
	if c != nil {
		return &IsometricNode{
			Node:  node,
			RGBA:  c,
			Order: node.Pos.Y() + ((bounds[1].X() - node.Pos.X()) * bounds[1].X()) + ((bounds[1].Z() - node.Pos.Z()) + bounds[1].Z()),
		}, nil
	}

	return nil, nil
}

func (r *IsoRenderer) getImagePos(x, y, z float64, size_x, size_y int) (float64, float64) {
	// max size of z or x axis
	max_xz := x
	if z > max_xz {
		max_xz = z
	}

	xpos := ((r.cubesize * x) - (r.cubesize * z)) + (float64(size_x) / 2) - r.cubesize
	ypos := float64(size_y) - (r.cubesize * sqrt3div2 * 2) -
		(r.cubesize * tan30 * x) -
		(r.cubesize * tan30 * z) -
		(r.cubesize * sqrt3div2 * y)

	return xpos, ypos
}
