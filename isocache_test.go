package maprenderer_test

import (
	"image/color"
	"testing"

	"github.com/minetest-go/maprenderer"
	"github.com/stretchr/testify/assert"
)

func TestGetCachedIsoCubeImage(t *testing.T) {

	img := maprenderer.GetCachedIsoCubeImage(&color.RGBA{R: 200, G: 100, B: 50, A: 0}, 20)
	assert.NotNil(t, img)
}
