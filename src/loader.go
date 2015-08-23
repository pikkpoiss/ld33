// Copyright 2015 Pikkpoiss
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"../lib/twodee"
	"fmt"
	"github.com/pikkpoiss/tmxgo"
	"io/ioutil"
)

func LoadTiledMap(path string) (grid *twodee.Grid, err error) {
	var (
		data  []byte
		m     *tmxgo.Map
		tiles []*tmxgo.Tile
		tile  *tmxgo.Tile
		x     int32
		y     int32
	)
	if data, err = ioutil.ReadFile(path); err != nil {
		return
	}
	if m, err = tmxgo.ParseMapString(string(data)); err != nil {
		return
	}
	if tiles, err = m.TilesFromLayerName("ground"); err != nil {
		return
	}
	grid = twodee.NewGrid(m.Width, m.Height, 1.0)
	for x = 0; x < grid.Width; x++ {
		for y = 0; y < grid.Height; y++ {
			tile = tiles[y*grid.Width+x]
			if tile != nil {
				grid.Set(
					x,
					y,
					&GridItem{
						passable: true,
						distance: -1,
						Frame:    fmt.Sprintf("tiles_%02v", tile.Index),
					},
				)
			}
		}
	}
	return
}
