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
	gridRenderer   *GridRenderer
	spriteSheet    *twodee.Spritesheet
	spriteTexture  *twodee.Texture
	level          *Level
	mouseX, mouseY float32
}

func NewGameLayer() (layer *GameLayer, err error) {
	var (
		gridRenderer *GridRenderer
		level        *Level
	)
	if level, err = NewLevel(); err != nil {
		return
	}
	layer = &GameLayer{
		level:        level,
		gridRenderer: gridRenderer,
	}
	err = layer.Reset()
	return
}

func (l *GameLayer) Delete() {
	l.gridRenderer.Delete()
}

func (l *GameLayer) Render() {
	l.spriteTexture.Bind()
	l.gridRenderer.Draw(l.level.Grid)
	l.spriteTexture.Unbind()
}

func (l *GameLayer) HandleEvent(evt twodee.Event) bool {
	switch event := evt.(type) {
	case *twodee.MouseMoveEvent:
		l.mouseX, l.mouseY = event.X, event.Y
	}
	return true
}

func (l *GameLayer) Reset() (err error) {
	var (
		camera *twodee.Camera
	)
	camera, err = twodee.NewCamera(
		twodee.Rect(0, 0, float32(l.level.Grid.Width()), float32(l.level.Grid.Height())),
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
