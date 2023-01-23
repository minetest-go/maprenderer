package maprenderer_test

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/minetest-go/mapparser"
	"github.com/minetest-go/maprenderer"
)

func NewMap() *Map {
	return &Map{
		hex_data: make(map[int64]string),
		map_data: make(map[int64]*mapparser.MapBlock),
	}
}

type Map struct {
	hex_data map[int64]string
	map_data map[int64]*mapparser.MapBlock
}

func (m Map) Load(csvfile string) error {
	file, err := os.Open(csvfile)
	if err != nil {
		return err
	}
	sc := bufio.NewScanner(file)
	line_num := 0
	for sc.Scan() {
		line := sc.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 2 {
			return fmt.Errorf("invalid format @ %d", line_num)
		}

		if len(parts[1])%2 != 0 {
			return fmt.Errorf("invalid hex count @ %d, len: %d", line_num, len(parts[1]))
		}

		pos, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return err
		}

		m.hex_data[pos] = parts[1]
		line_num++
	}

	return nil
}

func (m *Map) GetNode(pos *maprenderer.Pos) (*maprenderer.Node, error) {
	mb_pos := pos.Divide(16)
	pos_plain := CoordToPlain(mb_pos[0], mb_pos[1], mb_pos[2])
	if m.hex_data[pos_plain] == "" {
		// no mapblock there
		return nil, nil
	}

	md := m.map_data[pos_plain]
	if md == nil {
		b, err := hex.DecodeString(m.hex_data[pos_plain])
		if err != nil {
			return nil, err
		}

		md, err = mapparser.Parse(b)
		if err != nil {
			return nil, err
		}

		m.map_data[pos_plain] = md
	}

	rel_pos := pos.Subtract(mb_pos.Multiply(16))

	n := &maprenderer.Node{
		Name:   md.GetNodeName(rel_pos.X(), rel_pos.Y(), rel_pos.Z()),
		Param1: 0,
		Param2: md.GetParam2(rel_pos.X(), rel_pos.Y(), rel_pos.Z()),
		Pos:    pos,
	}

	return n, nil
}

func (m *Map) SearchNode(pos, direction *maprenderer.Pos, iterations int) (*maprenderer.Node, error) {
	current_pos := pos
	for i := 0; i <= iterations; i++ {
		node, err := m.GetNode(current_pos)
		if err != nil {
			return nil, err
		}

		if node != nil && node.Name != "air" {
			return node, nil
		}

		current_pos = current_pos.Add(direction)
	}

	return nil, nil
}

// utils

const (
	numBitsPerComponent = 12
	modulo              = 1 << numBitsPerComponent
	maxPositive         = modulo / 2
	minValue            = -1 << (numBitsPerComponent - 1)
	maxValue            = 1<<(numBitsPerComponent-1) - 1

	MinPlainCoord = -34351347711
)

func CoordToPlain(x, y, z int) int64 {
	return int64(z)<<(2*numBitsPerComponent) +
		int64(y)<<numBitsPerComponent +
		int64(x)
}

func unsignedToSigned(i int16) int {
	if i < maxPositive {
		return int(i)
	}
	return int(i - maxPositive*2)
}

// To match C++ code.
func pythonModulo(i int16) int16 {
	const mask = modulo - 1
	if i >= 0 {
		return i & mask
	}
	return modulo - -i&mask
}

func PlainToCoord(i int64) (int, int, int) {
	x := unsignedToSigned(pythonModulo(int16(i)))
	i = (i - int64(x)) >> numBitsPerComponent
	y := unsignedToSigned(pythonModulo(int16(i)))
	i = (i - int64(y)) >> numBitsPerComponent
	z := unsignedToSigned(pythonModulo(int16(i)))
	return x, y, z
}
