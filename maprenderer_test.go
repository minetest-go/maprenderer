package maprenderer

import (
	"fmt"
	"image/png"
	"os"
	"testing"

	"github.com/minetest-go/maprenderer/colormapping"
	"github.com/stretchr/testify/assert"
)

func TestMapRendererNotMultiple(t *testing.T) {
	r, err := NewMapRenderer(nil, nil, 11)
	assert.Nil(t, r)
	assert.Error(t, err)
}

func TestMapRenderer(t *testing.T) {
	m := NewMap()
	err := m.Load("testdata/map.csv")
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

	// test not lined up
	_, err = r.Render(&MapblockPos{X: 1, Y: 0, Z: 0}, &MapblockPos{X: 0, Y: 10, Z: 0})
	assert.Error(t, err)
	_, err = r.Render(&MapblockPos{X: 0, Y: 0, Z: 1}, &MapblockPos{X: 0, Y: 10, Z: 0})
	assert.Error(t, err)

	// error case in mapblock accessor
	_, err = r.Render(&MapblockPos{X: 666, Y: 0, Z: 0}, &MapblockPos{X: 0, Y: 10, Z: 0})
	assert.Error(t, err)

	for x := 0; x < 4; x++ {
		for z := 0; z < 4; z++ {
			pos1 := &MapblockPos{X: x, Y: 0, Z: z}
			pos2 := &MapblockPos{X: x, Y: 10, Z: z}
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
