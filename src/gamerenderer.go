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

type GameRenderer struct {
	sheet   *twodee.Spritesheet
	sprite  *twodee.SpriteRenderer
	effects *EffectsRenderer
}

func NewGameRenderer(level *Level, sheet *twodee.Spritesheet) (renderer *GameRenderer, err error) {
	var (
		xsize = int(PxPerUnit) * int(level.Grid.Width())
		ysize = int(PxPerUnit) * int(level.Grid.Height())
	)
	renderer = &GameRenderer{
		sheet: sheet,
	}
	if renderer.sprite, err = twodee.NewSpriteRenderer(level.Camera); err != nil {
		return
	}
	if renderer.effects, err = NewEffectsRenderer(xsize, ysize); err != nil {
		return
	}
	return
}

func (r *GameRenderer) Delete() {
	r.sprite.Delete()
	r.effects.Delete()
}

func (r *GameRenderer) Draw(level *Level) {
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
		configs = append(configs, mob.SpriteConfig(r.sheet))
	}
	configs = append(configs, r.cursorSpriteConfig(r.sheet, level.GetMouse(), level.GetCursor()))
	r.effects.Bind()
	r.sprite.Draw(configs)
	r.effects.Unbind()
	r.effects.Draw()
}

func (r *GameRenderer) cursorSpriteConfig(sheet *twodee.Spritesheet, pt mgl32.Vec2, cursor string) twodee.SpriteConfig {
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

func (r *GameRenderer) gridSpriteConfig(sheet *twodee.Spritesheet, x, y float32, item *GridItem) twodee.SpriteConfig {
	var frame *twodee.SpritesheetFrame
	if !item.Passable() {
		frame = sheet.GetFrame("special_squares_00")
	} else if item.Distance() >= 0 && item.Distance() < 15 {
		frame = sheet.GetFrame(fmt.Sprintf("numbered_squares_%02v", item.Distance()))
	} else {
		frame = sheet.GetFrame(item.Frame)
	}
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			x + frame.Width/2.0, y + frame.Height/2.0, 0,
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}
}
