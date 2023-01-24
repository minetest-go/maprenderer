package maprenderer

import (
	"fmt"
	"math"
)

type Pos [3]int

func NewPos(x, y, z int) *Pos {
	return &Pos{x, y, z}
}

func (p *Pos) X() int { return p[0] }
func (p *Pos) Y() int { return p[1] }
func (p *Pos) Z() int { return p[2] }

func (p *Pos) String() string {
	return fmt.Sprintf("Pos{%d,%d,%d}", p.X(), p.Y(), p.Z())
}

func (p1 *Pos) Add(p2 *Pos) *Pos {
	return &Pos{
		p1[0] + p2[0],
		p1[1] + p2[1],
		p1[2] + p2[2],
	}
}

func (p1 *Pos) Subtract(p2 *Pos) *Pos {
	return &Pos{
		p1[0] - p2[0],
		p1[1] - p2[1],
		p1[2] - p2[2],
	}
}

func (p *Pos) Divide(n float64) *Pos {
	return &Pos{
		int(math.Floor(float64(p[0]) / n)),
		int(math.Floor(float64(p[1]) / n)),
		int(math.Floor(float64(p[2]) / n)),
	}
}

func (p1 *Pos) Multiply(n int) *Pos {
	return &Pos{
		p1[0] * n,
		p1[1] * n,
		p1[2] * n,
	}
}
