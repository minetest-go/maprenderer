package maprenderer_test

import (
	"fmt"
	"image/png"
	"os"
	"testing"

	"github.com/minetest-go/colormapping"
	"github.com/minetest-go/maprenderer"
	"github.com/stretchr/testify/assert"
)

func TestMapRendererNotMultiple(t *testing.T) {
	r, err := maprenderer.NewMapRenderer(nil, nil, 11)
	assert.Nil(t, r)
	assert.Error(t, err)
}

func TestMapRenderer(t *testing.T) {
	m := NewMap()
	err := m.Load("testdata/map.csv")
	assert.NoError(t, err)

	cm := colormapping.NewColorMapping()
	assert.NotNil(t, cm)

	//load defaults
	err = cm.LoadDefaults()
	assert.NoError(t, err)

	r, err := maprenderer.NewMapRenderer(cm, m.GetMapblock, 16)
	assert.NoError(t, err)
	assert.NotNil(t, r)

	// error case in mapblock accessor
	_, err = r.Render(&maprenderer.MapblockPos{X: 666, Y: 0, Z: 0}, 10)
	assert.Error(t, err)

	for x := 0; x < 4; x++ {
		for z := 0; z < 4; z++ {
			pos1 := &maprenderer.MapblockPos{X: x, Y: 0, Z: z}
			img, err := r.Render(pos1, 10)
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

func ExampleMapRenderer() {
	// create color mapping
	cm := colormapping.NewColorMapping()

	// TODO: implemement mapblock fetching and parsing here
	accessor := func(pos maprenderer.MapblockPosGetter) (maprenderer.Mapblock, error) {
		return nil, nil
	}

	// create renderer with 256 px sidelength
	r, err := maprenderer.NewMapRenderer(cm, accessor, 256)
	if err != nil {
		panic(err)
	}

	// render the mapblock at 0,0,0 with 10 mapblocks y-height into an image
	img, err := r.Render(&maprenderer.MapblockPos{X: 0, Y: 0, Z: 0}, 10)
	if err != nil {
		panic(err)
	}

	// create an output file
	f, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}

	// write to the file
	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}
}
