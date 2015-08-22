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
	"io/ioutil"
	"time"
)

const (
	PxPerUnit = 32
)

type GameLayer struct {
	grid          *twodee.Grid
	gridRenderer  *GridRenderer
	spriteSheet   *twodee.Spritesheet
	spriteTexture *twodee.Texture
}

func NewGameLayer() (layer *GameLayer, err error) {
	var (
		gridRenderer *GridRenderer
		grid         *twodee.Grid
	)
	grid = twodee.NewGrid(64, 40, PxPerUnit)
	if err != nil {
		return
	}
	layer = &GameLayer{
		grid:         grid,
		gridRenderer: gridRenderer,
	}
	grid.Set(4, 19, &GridItem{false})
	grid.Set(60, 20, &GridItem{false})
	err = layer.Reset()
	return
}

func (l *GameLayer) Delete() {
	l.gridRenderer.Delete()
}

func (l *GameLayer) Render() {
	l.spriteTexture.Bind()
	l.gridRenderer.Draw(l.grid)
	l.spriteTexture.Unbind()
}

func (l *GameLayer) HandleEvent(evt twodee.Event) bool {
	return true
}

func (l *GameLayer) Reset() (err error) {
	var (
		camera *twodee.Camera
	)
	camera, err = twodee.NewCamera(
		twodee.Rect(0, 0, float32(l.grid.Width), float32(l.grid.Height)),
		twodee.Rect(0, 0, 1024, 640),
	)
	if err = l.loadSpritesheet(); err != nil {
		return
	}
	if l.gridRenderer, err = NewGridRenderer(camera, l.spriteSheet); err != nil {
		return
	}
	return
}

func (l *GameLayer) Update(elapsed time.Duration) {
}

func (l *GameLayer) loadSpritesheet() (err error) {
	var (
		data []byte
	)
	if data, err = ioutil.ReadFile("resources/spriteSheet.json"); err != nil {
		return
	}
	if l.spriteSheet, err = twodee.ParseTexturePackerJSONArrayString(
		string(data),
		PxPerUnit,
	); err != nil {
		return
	}
	if l.spriteTexture, err = twodee.LoadTexture(
		"resources/"+l.spriteSheet.TexturePath,
		twodee.Nearest,
	); err != nil {
		return
	}
	return
}
