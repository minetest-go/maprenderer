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

	ir, err := maprenderer.NewIsoRenderer(cm.GetColor, m, 5)
	assert.NoError(t, err)
	assert.NotNil(t, ir)

	from := &maprenderer.Pos{0, 16, 0}
	to := &maprenderer.Pos{16 - 1, 32 - 1, 16 - 1}
	img, err := ir.Render(from, to)
	assert.NoError(t, err)

	f, err := os.OpenFile("output/iso-test.png", os.O_CREATE|os.O_RDWR, 0755)
	assert.NoError(t, err)

	err = png.Encode(f, img)
	assert.NoError(t, err)
}

func BenchmarkIsoRenderMap(b *testing.B) {
	m := NewMap()
	err := m.Load("testdata/map.csv")
	assert.NoError(b, err)

	cm := colormapping.NewColorMapping()
	assert.NotNil(b, cm)

	//load defaults
	err = cm.LoadDefaults()
	assert.NoError(b, err)

	ir, err := maprenderer.NewIsoRenderer(cm.GetColor, m, 5)
	assert.NoError(b, err)
	assert.NotNil(b, ir)

	from := &maprenderer.Pos{0, 0, 0}
	to := &maprenderer.Pos{16 - 1, 16 - 1, 16 - 1}

	// prime cube image cache
	img, err := ir.Render(from, to)
	assert.NoError(b, err)
	assert.NotNil(b, img)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		img, err := ir.Render(from, to)
		assert.NoError(b, err)
		assert.NotNil(b, img)
	}
}
