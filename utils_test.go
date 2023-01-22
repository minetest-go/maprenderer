package maprenderer_test

import (
	"testing"

	"github.com/minetest-go/maprenderer"
	"github.com/stretchr/testify/assert"
)

func TestSortPos(t *testing.T) {
	p1 := [3]int{1, 2, 3}
	p2 := [3]int{3, 2, 1}

	p1, p2 = maprenderer.SortPos(p1, p2)
	assert.Equal(t, 1, p1[0])
	assert.Equal(t, 2, p1[1])
	assert.Equal(t, 1, p1[2])
	assert.Equal(t, 3, p2[0])
	assert.Equal(t, 2, p2[1])
	assert.Equal(t, 3, p2[2])
}

func TestNodePosToMapblock(t *testing.T) {
	assert.Equal(t, [3]int{0, 0, 0}, maprenderer.NodePosToMapblock([3]int{0, 0, 1}))
	assert.Equal(t, [3]int{0, 0, 0}, maprenderer.NodePosToMapblock([3]int{0, 15, 1}))
	assert.Equal(t, [3]int{0, 1, 0}, maprenderer.NodePosToMapblock([3]int{0, 16, 1}))
	assert.Equal(t, [3]int{0, 0, 0}, maprenderer.NodePosToMapblock([3]int{0, 0, 15}))
	assert.Equal(t, [3]int{0, 0, -1}, maprenderer.NodePosToMapblock([3]int{0, 0, -10}))
}
