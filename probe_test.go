package maprenderer_test

import (
	"testing"

	"github.com/minetest-go/colormapping"
	"github.com/minetest-go/maprenderer"
	"github.com/minetest-go/types"
	"github.com/stretchr/testify/assert"
)

func TestProbe(t *testing.T) {
	m := NewMap()
	err := m.Load("testdata/map.csv")
	assert.NoError(t, err)

	cm := colormapping.NewColorMapping()
	assert.NotNil(t, cm)
	err = cm.LoadDefaults()
	assert.NoError(t, err)

	min := types.NewPos(0, 0, 0)
	max := types.NewPos(32, 32, 32)
	pos := types.NewPos(0, 32, 0)
	ipos := types.NewPos(0, -1, 0)

	nodes, err := maprenderer.Probe(min, max, pos, ipos, m.GetNode, cm.GetColor, false)
	assert.NoError(t, err)
	assert.NotNil(t, nodes)
	//TODO
}
