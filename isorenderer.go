package maprenderer

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"sort"
)

var tan30 = math.Tan(30 * math.Pi / 180)
var sqrt3div2 = 2 / math.Sqrt(3)

func NewIsoRenderer(cr ColorResolver, na NodeAccessor, cubesize int) (*IsoRenderer, error) {
	return &IsoRenderer{
		cr:       cr,
		na:       na,
		cubesize: float64(cubesize),
	}, nil
}

type IsoRenderer struct {
	cr       ColorResolver
	na       NodeAccessor
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

	for y := from.Y(); y <= to.Y(); y++ {
		// right side
		for x := to.X() - 1; x >= from.X(); x-- {
			n, err := r.searchNode(&Pos{x, y, from.Z()}, direction, from, [2]*Pos{from, to})
			if err != nil {
				return nil, err
			}
			if n != nil {
				nodes = append(nodes, n)
			}
		}

		// left side
		for z := to.Z() - 1; z >= from.Z(); z-- {
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
	size_x, size_y := GetIsometricImageSize(from, to, int(r.cubesize))
	img := image.NewRGBA(image.Rect(0, 0, size_x, size_y))

	for _, node := range nodes {
		x, y := r.getImagePos(float64(node.Pos[0]), float64(node.Pos[1]), float64(node.Pos[2]), size_x, size_y)

		cube_img := GetCachedIsoCubeImage(node.RGBA, r.cubesize)
		p1 := image.Point{X: int(math.Floor(x)), Y: int(math.Floor(y))}
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
	xpos := ((r.cubesize * x) - (r.cubesize * z)) + (float64(size_x) / 2)
	ypos := 260 - (r.cubesize * tan30 * x) - (r.cubesize * tan30 * z) - (r.cubesize * sqrt3div2 * y)

	return xpos, ypos
}

func GetIsometricImageSize(from, to *Pos, cubesize int) (int, int) {
	cx_s, cy_s := GetIsoCubeSize(float64(cubesize))

	// max size of z or x axis
	x_diff := to.X() - from.X()
	y_diff := to.Y() - from.Y()
	z_diff := to.Z() - from.Z()
	max_xz := x_diff
	if z_diff > x_diff {
		max_xz = z_diff
	}

	size_x := int(cx_s) * (x_diff + z_diff) / 2
	size_y := int(cy_s) * (y_diff + max_xz) / 2

	return size_x, size_y
}
