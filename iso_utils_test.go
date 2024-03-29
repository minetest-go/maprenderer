package maprenderer_test

import (
	"testing"

	"github.com/minetest-go/maprenderer"
	"github.com/stretchr/testify/assert"
)

func TestGetIsoCubeSize(t *testing.T) {
	size_x, size_y := maprenderer.GetIsoCubeSize(1.0)
	assert.InDelta(t, 1.732, size_x, 0.1)
	assert.InDelta(t, 2.0, size_y, 0.1)
}

func TestGetIsometricImageSize(t *testing.T) {
	cubesize := 5.0

	// event side length
	size := maprenderer.NewPos(16, 16, 16)
	size_x, size_y := maprenderer.GetIsometricImageSize(size, cubesize)
	assert.Equal(t, 139, size_x)
	assert.Equal(t, 160, size_y)

	// uneven side length
	size = maprenderer.NewPos(16, 16, 32)
	size_x, size_y = maprenderer.GetIsometricImageSize(size, cubesize)
	assert.Equal(t, 208, size_x)
	assert.Equal(t, 240, size_y)
}

func TestGetImagePosEvenSides(t *testing.T) {
	cubesize := 5.0
	size := maprenderer.NewPos(16, 16, 16)
	size_x, size_y := maprenderer.GetIsometricImageSize(size, cubesize)

	// top/center node
	// TODO
	rel_pos := maprenderer.NewPos(15, 15, 15)
	x, y := maprenderer.GetImagePos(rel_pos, size, size_x, size_y, cubesize)
	assert.InDelta(t, 75, x, 0.1) // (size_x/2)-cubesize
	assert.InDelta(t, 0, y, 1)

	// bottom/center node
	rel_pos = maprenderer.NewPos(0, 0, 0)
	x, y = maprenderer.GetImagePos(rel_pos, size, size_x, size_y, cubesize)
	assert.InDelta(t, 75, x, 0.1)
	assert.InDelta(t, 172.5, y, 1) // size_y - cubesize_y
}
