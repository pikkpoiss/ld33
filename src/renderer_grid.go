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
	grid   *twodee.Grid
	sheet  *twodee.Spritesheet
	sprite *twodee.SpriteRenderer
}

func NewGridRenderer(grid *twodee.Grid, sheet *twodee.Spritesheet) (renderer *GridRenderer, err error) {
	var (
		camera *twodee.Camera
		sprite *twodee.SpriteRenderer
	)
	camera, err = twodee.NewCamera(twodee.Rect(0, 0, 50, 50), twodee.Rect(0, 0, 640, 480))
	if sprite, err = twodee.NewSpriteRenderer(camera); err != nil {
		return
	}
	renderer = &GridRenderer{
		grid:   grid,
		sprite: sprite,
		sheet:  sheet,
	}
	return
}

func (r *GridRenderer) Delete() {
	r.sprite.Delete()
}

func (b *GridRenderer) spriteConfig(sheet *twodee.Spritesheet) twodee.SpriteConfig {
	frame := sheet.GetFrame("numbered_squares_00")
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			1.0, 2.0, 0,
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}
}
