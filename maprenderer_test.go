package maprenderer

import (
	"bufio"
	"fmt"
	"image/png"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/minetest-go/mapparser"
	"github.com/minetest-go/maprenderer/colormapping"
	"github.com/stretchr/testify/assert"
)

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

type Map struct {
	world map[int64]string
}

func NewMap() *Map {
	return &Map{
		world: make(map[int64]string),
	}
}

func (m *Map) GetMapblock(pos *MapblockPos) (Mapblock, error) {
	pos_plain := CoordToPlain(pos.X, pos.Y, pos.Z)
	str := m.world[pos_plain]
	if str == "" {
		return nil, nil
	}

	b := make([]byte, len(str)/2)
	for i := 0; i < len(str); i += 2 {
		num, err := strconv.ParseUint(str[i:i+2], 16, 32)
		if err != nil {
			panic(err)
		}
		b[i/2] = byte(num)
	}

	mb, err := mapparser.Parse(b)
	if err != nil {
		panic(err)
	}

	return mb, nil
}

func (m *Map) Load() error {
	file, err := os.Open("testdata/map.csv")
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

		m.world[pos] = parts[1]
		line_num++
	}

	return nil
}

func TestMapRenderer(t *testing.T) {
	m := NewMap()
	err := m.Load()
	assert.NoError(t, err)

	cm := colormapping.NewColorMapping()
	assert.NotNil(t, cm)

	data, err := os.ReadFile("testdata/mtg.txt")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	//load testcolors
	_, err = cm.LoadColorMapping(data)
	assert.NoError(t, err)

	r, err := NewMapRenderer(cm, m.GetMapblock, 16)
	assert.NoError(t, err)
	assert.NotNil(t, r)

	for x := 0; x < 4; x++ {
		for z := 0; z < 4; z++ {
			pos1 := MapblockPos{X: x, Y: 0, Z: z}
			pos2 := MapblockPos{X: x, Y: 10, Z: z}
			img, err := r.Render(pos1, pos2)
			assert.NoError(t, err)
			assert.NotNil(t, img)

			os.Mkdir("output", 0755)

			f, err := os.Create(fmt.Sprintf("output/output-%d-%d.png", x, z))
			assert.NoError(t, err)
			assert.NotNil(t, f)

			err = png.Encode(f, img)
			assert.NoError(t, err)
		}
	}
}
