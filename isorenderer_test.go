package maprenderer

import (
	"os"
	"testing"

	"github.com/minetest-go/maprenderer/colormapping"
	"github.com/stretchr/testify/assert"
)

func TestIsoRenderer(t *testing.T) {
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

	r, err := NewIsoRenderer(cm, m.GetMapblock, 64)
	assert.NoError(t, err)
	assert.NotNil(t, r)

	/*
		pos1 := &MapblockPos{X: 0, Y: 0, Z: 0}
		pos2 := &MapblockPos{X: 0, Y: 10, Z: 0}
		img, err := r.Render(pos1, pos2, IsoDirectionNorthEast)
		assert.NotNil(t, img)
		assert.NoError(t, err)
	*/
}
