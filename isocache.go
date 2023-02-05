package maprenderer

import (
	"fmt"
	"image"
	"image/color"
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
		size_x, size_y := GetIsoCubeSize(size)
		dc := gg.NewContext(int(size_x), int(size_y))
		DrawCube(dc, c, size, 0, 0)
		img = dc.Image()

		// cache for future use
		rc.lock.Lock()
		rc.cache[key] = img
		rc.lock.Unlock()
	}

	return img
}
