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
)

type GridRenderer struct {
	sheet  *twodee.Spritesheet
	sprite *twodee.SpriteRenderer
}

func NewGridRenderer(camera *twodee.Camera, sheet *twodee.Spritesheet) (renderer *GridRenderer, err error) {
	var (
		sprite *twodee.SpriteRenderer
	)
	if sprite, err = twodee.NewSpriteRenderer(camera); err != nil {
		return
	}
	renderer = &GridRenderer{
		sprite: sprite,
		sheet:  sheet,
	}
	return
}

func (r *GridRenderer) Delete() {
	r.sprite.Delete()
}

func (r *GridRenderer) Draw(grid *twodee.Grid) {
	var (
		configs = []twodee.SpriteConfig{}
		x       int32
		y       int32
		tditem  twodee.GridItem
		item    *GridItem
		ok      bool
	)
	for x = 0; x < grid.Width; x++ {
		for y = 0; y < grid.Height; y++ {
			tditem = grid.Get(x, y)
			if item, ok = tditem.(*GridItem); !ok {
				item = nil
			}
			configs = append(configs, r.spriteConfig(
				r.sheet,
				int(x),
				int(y),
				item,
			))
		}
	}
	r.sprite.Draw(configs)
}

func (r *GridRenderer) spriteConfig(sheet *twodee.Spritesheet, x, y int, item *GridItem) twodee.SpriteConfig {
	var frame *twodee.SpritesheetFrame
	if item == nil || item.Passable() {
		frame = sheet.GetFrame("numbered_squares_00")
	} else {
		frame = sheet.GetFrame("numbered_squares_14")
	}
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			float32(x) + 0.5, float32(y) + 0.5, 0,
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}
}
