package maprenderer_test

import (
	"image/png"
	"os"
	"testing"

	"github.com/minetest-go/colormapping"
	"github.com/minetest-go/maprenderer"
	"github.com/stretchr/testify/assert"
)

func TestIsoRenderer(t *testing.T) {
	m := NewMap()
	err := m.Load("testdata/map.csv")
	assert.NoError(t, err)

	cm := colormapping.NewColorMapping()
	assert.NotNil(t, cm)

	//load defaults
	err = cm.LoadDefaults()
	assert.NoError(t, err)

	from := maprenderer.NewPos(0, -10, 0)
	to := maprenderer.NewPos(160, 60, 160)
	opts := &maprenderer.IsoRenderOpts{
		CubeLen: 8,
	}

	img, err := maprenderer.RenderIsometric(m.GetNode, cm.GetColor, from, to, opts)
	assert.NoError(t, err)
	assert.NotNil(t, img)

	f, err := os.OpenFile("output/iso-test.png", os.O_CREATE|os.O_RDWR, 0755)
	assert.NoError(t, err)

	err = png.Encode(f, img)
	assert.NoError(t, err)
}
