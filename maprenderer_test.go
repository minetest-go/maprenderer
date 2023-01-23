package maprenderer_test

import (
	"image/png"
	"os"
	"testing"

	"github.com/minetest-go/colormapping"
	"github.com/minetest-go/maprenderer"
	"github.com/stretchr/testify/assert"
)

func TestRenderMap(t *testing.T) {
	m := NewMap()
	err := m.Load("testdata/map.csv")
	assert.NoError(t, err)

	cm := colormapping.NewColorMapping()
	assert.NotNil(t, cm)

	//load defaults
	err = cm.LoadDefaults()
	assert.NoError(t, err)

	from := &maprenderer.Pos{0, 0, 0}
	to := &maprenderer.Pos{128 - 1, 47, 128 - 1}

	img, err := maprenderer.RenderMap(from, to, m, cm.GetColor)
	assert.NoError(t, err)
	assert.NotNil(t, img)

	f, err := os.OpenFile("output/map-test.png", os.O_CREATE|os.O_RDWR, 0755)
	assert.NoError(t, err)

	err = png.Encode(f, img)
	assert.NoError(t, err)
}

func BenchmarkRenderMap(b *testing.B) {
	m := NewMap()
	err := m.Load("testdata/map.csv")
	assert.NoError(b, err)

	cm := colormapping.NewColorMapping()
	assert.NotNil(b, cm)

	//load defaults
	err = cm.LoadDefaults()
	assert.NoError(b, err)

	from := &maprenderer.Pos{0, 0, 0}
	to := &maprenderer.Pos{16 - 1, 16 - 1, 16 - 1}

	for i := 0; i < b.N; i++ {
		img, err := maprenderer.RenderMap(from, to, m, cm.GetColor)
		assert.NoError(b, err)
		assert.NotNil(b, img)
	}
}
