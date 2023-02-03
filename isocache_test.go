package maprenderer_test

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"testing"

	"github.com/minetest-go/maprenderer"
	"github.com/stretchr/testify/assert"
)

func TestGetCachedIsoCubeImage(t *testing.T) {

	size := 20.0

	size_x, size_y := maprenderer.GetIsoCubeSize(size)

	rc := maprenderer.NewIsoRenderCache()
	cube := rc.GetCachedIsoCubeImage(&color.RGBA{R: 200, G: 100, B: 50, A: 255}, size)
	assert.NotNil(t, cube)

	img := image.NewRGBA(image.Rect(0, 0, int(size_x), int(size_y)))

	p1 := image.Point{X: 0, Y: 0}
	draw.Draw(img, image.Rectangle{p1, p1.Add(cube.Bounds().Size())}, cube, image.Point{0, 0}, draw.Src)

	f, err := os.OpenFile("output/isocache-test.png", os.O_CREATE|os.O_RDWR, 0755)
	assert.NoError(t, err)

	err = png.Encode(f, img)
	assert.NoError(t, err)
}
