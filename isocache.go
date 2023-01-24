package maprenderer

import (
	"fmt"
	"image"
	"image/color"
	"sync"

	"github.com/fogleman/gg"
)

var cubeCache = make(map[string]image.Image)
var cubeCacheLock = sync.RWMutex{}

func getRGBASizeKey(c *color.RGBA, size float64) string {
	return fmt.Sprintf("%d/%d/%d/%d/%f", c.R, c.G, c.B, c.A, size)
}

func GetCachedIsoCubeImage(c *color.RGBA, size float64) image.Image {
	key := getRGBASizeKey(c, size)

	cubeCacheLock.RLock()
	img := cubeCache[key]
	cubeCacheLock.RUnlock()

	x := 10.0
	y := 10.0

	if img == nil {
		// create image

		dc := gg.NewContext(100, 100)

		// right side
		dc.MoveTo(size+x, (size*tan30)+y)
		dc.LineTo(x, (size*sqrt3div2)+y)
		dc.LineTo(x, y)
		dc.LineTo(size+x, -(size*tan30)+y)
		dc.ClosePath()
		dc.SetRGB255(int(c.R), int(c.G), int(c.B))
		dc.Fill()

		// left side
		dc.MoveTo(x, (size*sqrt3div2)+y)
		dc.LineTo(-size+x, (size*tan30)+y)
		dc.LineTo(-size+x, -(size*tan30)+y)
		dc.LineTo(x, y)
		dc.ClosePath()
		AdjustAndFill(dc, int(c.R), int(c.G), int(c.B), -20)
		dc.Fill()

		// top side
		dc.MoveTo(-size+x, -(size*tan30)+y)
		dc.LineTo(x, -(size*sqrt3div2)+y)
		dc.LineTo(size+x, -(size*tan30)+y)
		dc.LineTo(x, y)
		dc.ClosePath()
		AdjustAndFill(dc, int(c.R), int(c.G), int(c.B), 20)
		dc.Fill()

		img = dc.Image()

		// cache for future use
		cubeCacheLock.Lock()
		cubeCache[key] = img
		cubeCacheLock.Unlock()
	}

	return img
}
