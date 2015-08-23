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
	"io/ioutil"
)

type SplashRenderer struct {
	sheet   *twodee.Spritesheet
	sprite  *twodee.SpriteRenderer
	texture *twodee.Texture
}

func NewSplashRenderer(camera *twodee.Camera) (renderer *SplashRenderer, err error) {
	renderer = &SplashRenderer{}
	if err = renderer.loadSpritesheet(); err != nil {
		return
	}
	if renderer.sprite, err = twodee.NewSpriteRenderer(camera); err != nil {
		return
	}
	return
}

func (r *SplashRenderer) Delete() {
	r.sprite.Delete()
	r.texture.Delete()
}

func (r *SplashRenderer) Draw(state *State) {
	var (
		configs []twodee.SpriteConfig
		frame   string = ""
	)
	switch state.SplashState {
	case SplashStart:
		frame = "start.fw"
	case SplashWin:
		frame = "win.fw"
	case SplashLose:
		frame = "lose.fw"
	}
	if frame != "" {
		configs = []twodee.SpriteConfig{
			r.splashConfig(r.sheet, mgl32.Vec2{0, 0}, frame),
		}
		r.texture.Bind()
		r.sprite.Draw(configs)
		r.texture.Unbind()
	}
}

func (r *SplashRenderer) loadSpritesheet() (err error) {
	var (
		data []byte
	)
	if data, err = ioutil.ReadFile("resources/splash.json"); err != nil {
		return
	}
	if r.sheet, err = twodee.ParseTexturePackerJSONArrayString(
		string(data),
		PxPerUnit,
	); err != nil {
		return
	}
	if r.texture, err = twodee.LoadTexture(
		"resources/"+r.sheet.TexturePath,
		twodee.Nearest,
	); err != nil {
		return
	}
	return
}

func (r *SplashRenderer) splashConfig(sheet *twodee.Spritesheet, pt mgl32.Vec2, name string) twodee.SpriteConfig {
	frame := sheet.GetFrame(name)
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			pt.X() + frame.Width/2.0, pt.Y() + frame.Height/2.0, 0.0, // Center
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}
}
