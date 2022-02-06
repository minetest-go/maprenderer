package colormapping

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMapping(t *testing.T) {
	m := NewColorMapping()
	assert.NotNil(t, m)

	data, err := os.ReadFile("testdata/testcolor.txt")
	assert.NoError(t, err)
	assert.NotNil(t, data)

	// empty mapping
	assert.Equal(t, 0, len(m.colors))

	// no colors loaded
	c := m.GetColor("default:river_water_flowing", 0)
	assert.Nil(t, c)

	//load testcolors
	count, err := m.LoadColorMapping(data)
	assert.NoError(t, err)
	assert.Equal(t, 2, count)

	// test colors loaded
	c = m.GetColor("default:river_water_flowing", 0)
	assert.NotNil(t, c)
	assert.Equal(t, uint8(0), c.R)
	assert.Equal(t, uint8(0), c.G)
	assert.Equal(t, uint8(100), c.B)
	assert.Equal(t, uint8(255), c.A)
	c = m.GetColor("scifi_nodes:blacktile2", 0)
	assert.NotNil(t, c)
	assert.Equal(t, uint8(10), c.R)
	assert.Equal(t, uint8(20), c.G)
	assert.Equal(t, uint8(30), c.B)
	assert.Equal(t, uint8(255), c.A)

	// load param2 palette
	palette_data, err := os.ReadFile("testdata/unifieddyes_palette_extended.png")
	assert.NoError(t, err)
	assert.NotNil(t, data)
	palette_nodes, err := os.ReadFile("testdata/palette_nodes.txt")
	assert.NoError(t, err)
	assert.NotNil(t, data)
	err = m.LoadPalette(palette_data, palette_nodes)
	assert.NoError(t, err)

	// test node with extended palette
	c = m.GetColor("unifiedbricks:brickblock_multicolor_dark", 50)
	assert.NotNil(t, c)
	assert.Equal(t, uint8(255), c.R)
	assert.Equal(t, uint8(191), c.G)
	assert.Equal(t, uint8(129), c.B)
	assert.Equal(t, uint8(255), c.A)
	c = m.GetColor("unifiedbricks:brickblock_multicolor_dark", 20)
	assert.NotNil(t, c)
	assert.Equal(t, uint8(255), c.R)
	assert.Equal(t, uint8(191), c.G)
	assert.Equal(t, uint8(255), c.B)
	assert.Equal(t, uint8(255), c.A)
}
