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
	"github.com/go-gl/mathgl/mgl32"
)

type GridRenderer struct {
	sheet  *twodee.Spritesheet
	sprite *twodee.SpriteRenderer
}

func NewGridRenderer(level *Level, sheet *twodee.Spritesheet) (renderer *GridRenderer, err error) {
	var (
		sprite *twodee.SpriteRenderer
	)
	if sprite, err = twodee.NewSpriteRenderer(level.Camera); err != nil {
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

func (r *GridRenderer) Draw(level *Level) {
	var (
		configs = []twodee.SpriteConfig{}
		x       int32
		y       int32
		item    *GridItem
	)
	for x = 0; x < level.Grid.Width(); x++ {
		for y = 0; y < level.Grid.Height(); y++ {
			item = level.Grid.Get(Ivec2{x, y})
			configs = append(configs, r.gridSpriteConfig(
				r.sheet,
				float32(x),
				float32(y),
				item,
			))
		}
	}
	for _, mob := range level.Mobs {
		if !mob.Enabled { // No enabled mobs after first disabled mob.
			break
		}
		configs = append(configs, r.mobSpriteConfig(
			r.sheet,
			mob.Pos,
			mob,
		))
	}
	configs = append(configs, r.cursorSpriteConfig(r.sheet, level.GetMouse(), level.GetCursor()))
	r.sprite.Draw(configs)
}

func (r *GridRenderer) cursorSpriteConfig(sheet *twodee.Spritesheet, pt mgl32.Vec2, cursor string) twodee.SpriteConfig {
	frame := sheet.GetFrame(cursor)
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			pt.X(), pt.Y(), 0,
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}
}

func (r *GridRenderer) mobSpriteConfig(sheet *twodee.Spritesheet, pt mgl32.Vec2, mob *Mob) twodee.SpriteConfig {
	frame := sheet.GetFrame("numbered_squares_00")
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			pt.X(), pt.Y(), 0,
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}
}

func (r *GridRenderer) gridSpriteConfig(sheet *twodee.Spritesheet, x, y float32, item *GridItem) twodee.SpriteConfig {
	var frame *twodee.SpritesheetFrame
	if item.Distance() >= 0 && item.Distance() < 15 {
		frame = sheet.GetFrame(fmt.Sprintf("numbered_squares_%02v", item.Distance()))
	} else if item.Passable() {
		frame = sheet.GetFrame("numbered_squares_00")
	} else {
		frame = sheet.GetFrame("numbered_squares_14")
	}
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			x + 0.5, y + 0.5, 0,
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}
}
