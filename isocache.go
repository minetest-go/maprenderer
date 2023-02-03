package maprenderer

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"sync"

	"github.com/fogleman/gg"
)

type IsoRenderCache struct {
	cache map[string]image.Image
	lock  sync.RWMutex
}

func NewIsoRenderCache() *IsoRenderCache {
	return &IsoRenderCache{
		cache: make(map[string]image.Image),
		lock:  sync.RWMutex{},
	}
}

func (rc *IsoRenderCache) getRGBASizeKey(c *color.RGBA, size float64) string {
	return fmt.Sprintf("%d/%d/%d/%d/%f", c.R, c.G, c.B, c.A, size)
}

func (rc *IsoRenderCache) GetCachedIsoCubeImage(c *color.RGBA, size float64) image.Image {
	key := rc.getRGBASizeKey(c, size)

	rc.lock.RLock()
	img := rc.cache[key]
	rc.lock.RUnlock()

	if img == nil {
		// create image
		size_x, size_y := GetIsoCubeSize(size)
		// round up
		size_x = math.Ceil(size_x)
		size_y = math.Ceil(size_y)

		// center position
		center_x := size_x / 2
		center_y := size_y / 2

		// proportional size
		sin30_proportional := sin30 * size

		dc := gg.NewContext(int(math.Ceil(size_x)), int(math.Ceil(size_y)))

		// right side
		dc.SetRGBA255(int(c.R), int(c.G), int(c.B), int(c.A))
		dc.MoveTo(center_x, center_y)
		dc.LineTo(size_x, center_y-sin30_proportional)
		dc.LineTo(size_x, size_y-sin30_proportional)
		dc.LineTo(center_x, size_y)
		dc.ClosePath()
		dc.Fill()

		// left side
		dc.SetRGBA255(
			AdjustColorComponent(c.R, -20),
			AdjustColorComponent(c.G, -20),
			AdjustColorComponent(c.B, -20),
			int(c.A),
		)
		dc.MoveTo(center_x, center_y)
		dc.LineTo(center_x, size_y)
		dc.LineTo(0, size_y-sin30_proportional)
		dc.LineTo(0, center_y-sin30_proportional)
		dc.ClosePath()
		dc.Fill()

		// top side
		dc.SetRGBA255(
			AdjustColorComponent(c.R, 20),
			AdjustColorComponent(c.G, 20),
			AdjustColorComponent(c.B, 20),
			int(c.A),
		)
		dc.MoveTo(center_x, center_y)
		dc.LineTo(0, center_y-sin30_proportional)
		dc.LineTo(center_x, 0)
		dc.LineTo(size_x, center_y-sin30_proportional)
		dc.ClosePath()
		dc.Fill()

		img = dc.Image()

		// cache for future use
		rc.lock.Lock()
		rc.cache[key] = img
		rc.lock.Unlock()
	}

	return img
}
