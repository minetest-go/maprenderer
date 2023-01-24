package maprenderer

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/fogleman/gg"
)

var cubeCache = make(map[string]image.Image)
var cubeCacheLock = sync.RWMutex{}

func getRGBASizeKey(c *color.RGBA, size float64) string {
	return fmt.Sprintf("%d/%d/%d/%d/%f", c.R, c.G, c.B, c.A, size)
}

func GetIsoCubeSize(size float64) (float64, float64) {
	return size * 2, size * sqrt3div2 * 2
}

func GetCachedIsoCubeImage(c *color.RGBA, size float64) image.Image {
	key := getRGBASizeKey(c, size)

	cubeCacheLock.RLock()
	img := cubeCache[key]
	cubeCacheLock.RUnlock()

	if img == nil {
		// create image
		size_x, size_y := GetIsoCubeSize(size)
		// center position
		x := size_x / 2
		y := size_y / 2

		dc := gg.NewContext(int(math.Ceil(size_x)), int(math.Ceil(size_y)))

		// right side
		dc.SetRGBA255(int(c.R), int(c.G), int(c.B), int(c.A))
		dc.MoveTo(size+x, (size*tan30)+y)
		dc.LineTo(x, (size*sqrt3div2)+y)
		dc.LineTo(x, y)
		dc.LineTo(size+x, -(size*tan30)+y)
		dc.ClosePath()
		dc.FillPreserve()
		dc.Stroke()

		// left side
		dc.SetRGBA255(
			AdjustColorComponent(c.R, -20),
			AdjustColorComponent(c.G, -20),
			AdjustColorComponent(c.B, -20),
			int(c.A),
		)
		dc.MoveTo(x, (size*sqrt3div2)+y)
		dc.LineTo(-size+x, (size*tan30)+y)
		dc.LineTo(-size+x, -(size*tan30)+y)
		dc.LineTo(x, y)
		dc.ClosePath()
		dc.FillPreserve()
		dc.Stroke()

		// top side
		dc.SetRGBA255(
			AdjustColorComponent(c.R, 20),
			AdjustColorComponent(c.G, 20),
			AdjustColorComponent(c.B, 20),
			int(c.A),
		)
		dc.MoveTo(-size+x, -(size*tan30)+y)
		dc.LineTo(x, -(size*sqrt3div2)+y)
		dc.LineTo(size+x, -(size*tan30)+y)
		dc.LineTo(x, y)
		dc.ClosePath()
		dc.FillPreserve()
		dc.Stroke()

		img = dc.Image()

		// cache for future use
		cubeCacheLock.Lock()
		cubeCache[key] = img
		cubeCacheLock.Unlock()
	}

	return img
}
