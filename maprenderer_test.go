package maprenderer_test

import (
	"image/png"
	"os"
	"testing"

	"github.com/minetest-go/colormapping"
	"github.com/minetest-go/maprenderer"
	"github.com/minetest-go/types"
	"github.com/stretchr/testify/assert"
)

func TestRenderMap(t *testing.T) {
	// map
	m := NewMap()
	err := m.Load("testdata/map.csv")
	assert.NoError(t, err)

	// colormapping
	cm := colormapping.NewColorMapping()
	assert.NotNil(t, cm)
	err = cm.LoadDefaults()
	assert.NoError(t, err)

	from := types.NewPos(0, -20, 0)
	to := types.NewPos(100, 50, 100)
	opts := &maprenderer.MapRenderOpts{}

	img, err := maprenderer.RenderMap(m.GetNode, cm.GetColor, from, to, opts)
	assert.NoError(t, err)
	assert.NotNil(t, img)

	f, err := os.OpenFile("output/map-test.png", os.O_CREATE|os.O_RDWR, 0755)
	assert.NoError(t, err)

	err = png.Encode(f, img)
	assert.NoError(t, err)
}
