package colormapping

import (
	"bufio"
	"bytes"
	"errors"
	"image/color"
	"io"
	"io/ioutil"
	"strconv"
	"strings"
)

type ColorMapping struct {
	colors               map[string]*color.RGBA
	extendedpaletteblock map[string]bool
	extendedpalette      *Palette
}

func (m *ColorMapping) GetColor(name string, param2 int) *color.RGBA {
	//TODO: list of node->palette
	if m.extendedpaletteblock[name] {
		// param2 coloring
		return m.extendedpalette.GetColor(param2)
	}

	return m.colors[name]
}

func (m *ColorMapping) GetColors() map[string]*color.RGBA {
	return m.colors
}

func (m *ColorMapping) ReadColorMapping(r io.Reader) (int, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return 0, err
	}

	return m.LoadColorMapping(b)
}

func (m *ColorMapping) LoadColorMapping(buffer []byte) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(buffer))
	count := 0
	line := 0

	for scanner.Scan() {
		line++

		txt := strings.Trim(scanner.Text(), " ")

		if len(txt) == 0 {
			//empty
			continue
		}

		if strings.HasPrefix(txt, "#") {
			//comment
			continue
		}

		parts := strings.Fields(txt)

		if len(parts) < 4 {
			return 0, errors.New("invalid line: #" + strconv.Itoa(line))
		}

		if len(parts) >= 4 {
			r, err := strconv.ParseInt(parts[1], 10, 32)
			if err != nil {
				return 0, err
			}

			g, err := strconv.ParseInt(parts[2], 10, 32)
			if err != nil {
				return 0, err
			}

			b, err := strconv.ParseInt(parts[3], 10, 32)
			if err != nil {
				return 0, err
			}

			a := int64(255)

			/*
				if len(parts) >= 5 {
					//with alpha
					//a, err = strconv.ParseInt(parts[4], 10, 32)
					//if err != nil {
					//	return 0, err
					//}
				}
			*/

			c := color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
			m.colors[parts[0]] = &c
			count++
		}
	}

	return count, nil
}

func (m *ColorMapping) LoadPalette(palettefile, nodelistfile []byte) error {
	palette, err := NewPalette(palettefile)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewReader(nodelistfile))
	paletteblock := make(map[string]bool)

	if err != nil {
		return err
	}

	for scanner.Scan() {
		txt := strings.Trim(scanner.Text(), " ")

		if len(txt) == 0 {
			//empty
			continue
		}

		paletteblock[txt] = true
	}

	m.extendedpaletteblock = paletteblock
	m.extendedpalette = palette
	return nil
}

func NewColorMapping() *ColorMapping {
	return &ColorMapping{
		colors: make(map[string]*color.RGBA),
	}
}
