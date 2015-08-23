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
	"github.com/go-gl/mathgl/mgl32"
	"sort"
)

type ByY []twodee.SpriteConfig

func (a ByY) Len() int      { return len(a) }
func (a ByY) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByY) Less(i, j int) bool {
	if a[i].View.Y == a[j].View.Y {
		return a[i].View.X < a[j].View.X
	}
	return a[i].View.Y > a[j].View.Y
}

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
		count      = level.Grid.Width() * level.Grid.Height()
		configs    = make([]twodee.SpriteConfig, 0, count)
		bgs        = make([]twodee.SpriteConfig, 0, count)
		highlights = make([]twodee.SpriteConfig, 0, len(level.Highlights))
		x          int32
		y          int32
		item       *GridItem
		pt         Ivec2
	)
	for x = 0; x < level.Grid.Width(); x++ {
		for y = 0; y < level.Grid.Height(); y++ {
			pt = Ivec2{x, y}
			item = level.Grid.GetBg(pt)
			bgs = append(bgs, r.gridSpriteConfig(
				r.sheet,
				float32(x),
				float32(y),
				item,
			))
			item = level.Grid.Get(pt)
			if item == nil {
				continue
			}
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
	for _, highlight := range level.Highlights {
		highlights = append(highlights, r.highlightSpriteConfig(r.sheet, highlight.Pos, highlight.Frame))
	}
	configs = append(configs, r.cursorSpriteConfig(r.sheet, level.GetMouse(), level.GetCursor()))
	sort.Sort(ByY(configs))
	r.effects.Bind()
	r.sprite.Draw(bgs)
	if len(highlights) > 0 {
		r.sprite.Draw(highlights)
	}
	r.sprite.Draw(configs)
	r.effects.Unbind()
	r.effects.Draw()
}

func (r *GameRenderer) highlightSpriteConfig(sheet *twodee.Spritesheet, pt Ivec2, name string) twodee.SpriteConfig {
	frame := sheet.GetFrame(name)
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			float32(pt.X()) + frame.Width/2.0, float32(pt.Y()) + frame.Height/2.0, 0.0,
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}
}

func (r *GameRenderer) cursorSpriteConfig(sheet *twodee.Spritesheet, pt mgl32.Vec2, cursor string) twodee.SpriteConfig {
	frame := sheet.GetFrame(cursor)
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			pt.X(), pt.Y(), 0.2,
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}
}

func (r *GameRenderer) gridSpriteConfig(sheet *twodee.Spritesheet, x, y float32, item *GridItem) twodee.SpriteConfig {
	var frame *twodee.SpritesheetFrame
	frame = sheet.GetFrame(item.Frame)
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			x + frame.Width/2.0, y + frame.Height/2.0, 0.0,
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}
}
