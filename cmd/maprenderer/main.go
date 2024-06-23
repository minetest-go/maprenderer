package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/minetest-go/colormapping"
	"github.com/minetest-go/mapparser"
	"github.com/minetest-go/maprenderer"
	"github.com/minetest-go/mtdb"
	"github.com/minetest-go/mtdb/block"
	"github.com/minetest-go/types"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var help = flag.Bool("help", false, "shows the help")
var show_version = flag.Bool("version", false, "shows the version")
var map_type = flag.String("type", "map", "map-type to render: 'map' or 'isometric'")
var output = flag.String("output", "map.png", "png to write to")
var from_pos = flag.String("from", "0,0,0", "from-position to use, in 'x,y,z' format")
var to_pos = flag.String("to", "100,100,100", "to-position to use, in 'x,y,z' format")
var world_dir = flag.String("world", "./", "world directory")

func parsePos(str string) (*types.Pos, error) {
	parts := strings.Split(str, ",")
	if len(parts) != 3 {
		return nil, fmt.Errorf("unexpected part-count: %d", len(parts))
	}

	x, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("could not parse x-component: '%s'", parts[0])
	}

	y, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("could not parse y-component: '%s'", parts[0])
	}

	z, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("could not parse z-component: '%s'", parts[0])
	}

	return types.NewPos(int(x), int(y), int(z)), nil
}

func NewNodeAccessor(repo block.BlockRepository) types.NodeAccessor {
	cache := expirable.NewLRU[string, *types.MapBlock](100, nil, time.Second*10)

	return func(p *types.Pos) (*types.Node, error) {
		mb_pos := p.Divide(16)
		key := mb_pos.String()

		mb, found := cache.Get(key)
		if mb == nil && found {
			return nil, nil
		}

		if !found {
			block, err := repo.GetByPos(mb_pos.X(), mb_pos.Y(), mb_pos.Z())
			if err != nil {
				return nil, fmt.Errorf("GetByPos error @ %s: %v", p, err)
			}

			if block != nil {
				mb, err = mapparser.Parse(block.Data)
				if err != nil {
					return nil, fmt.Errorf("parse error @ %s: %v", p, err)
				}
				if !mb.AirOnly {
					cache.Add(key, mb)
				}
			} else {
				cache.Add(key, nil)
				return nil, nil
			}
		}

		rel_pos := p.Subtract(mb_pos.Multiply(16))
		return &types.Node{
			Name:   mb.GetNodeName(rel_pos),
			Param1: 0,
			Param2: mb.GetParam2(rel_pos),
			Pos:    p,
		}, nil
	}
}

func main() {
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	if *show_version {
		fmt.Printf("maprenderer %s, commit %s, built at %s\n", version, commit, date)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if *world_dir != "./" {
		wd = *world_dir
	}

	repo, err := mtdb.NewBlockDB(wd)
	if err != nil {
		panic(err)
	}
	na := NewNodeAccessor(repo)

	cm := colormapping.NewColorMapping()
	err = cm.LoadDefaults()
	if err != nil {
		panic(err)
	}

	from, err := parsePos(*from_pos)
	if err != nil {
		panic(err)
	}

	to, err := parsePos(*to_pos)
	if err != nil {
		panic(err)
	}

	var img image.Image
	if *map_type == "isometric" {
		opts := maprenderer.NewDefaultIsoRenderOpts()
		img, err = maprenderer.RenderIsometric(na, cm.GetColor, from, to, opts)
	} else {
		opts := maprenderer.NewDefaultMapRenderOpts()
		img, err = maprenderer.RenderMap(na, cm.GetColor, from, to, opts)
	}

	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(*output, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}

	err = png.Encode(f, img)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Map rendering complete\n")
}
