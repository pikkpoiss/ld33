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
	"image/color"
	"time"
	"strconv"
)

type HudLayer struct {
	camera       *twodee.Camera
	textRenderer *twodee.TextRenderer
	regfont      *twodee.FontFace
	state        *State
	app          *Application
}

func NewHudLayer(state *State, app *Application) (layer *HudLayer, err error) {
	var (
		textRenderer *twodee.TextRenderer
		regfont      *twodee.FontFace
		bg           = color.Transparent
		font         = "resources/fonts/Prototype.ttf"
		camera       *twodee.Camera
	)
	if camera, err = twodee.NewCamera(
		twodee.Rect(0, 0, ScreenWidth, ScreenHeight),
		twodee.Rect(0, 0, ScreenWidth, ScreenHeight),
	); err != nil {
		return
	}
	if textRenderer, err = twodee.NewTextRenderer(camera); err != nil {
		return
	}
	if regfont, err = twodee.NewFontFace(font, 32, color.RGBA{100, 100, 100, 255}, bg); err != nil {
		return
	}
	layer = &HudLayer{
		camera:       camera,
		textRenderer: textRenderer,
		regfont:      regfont,
		state:        state,
		app:          app,
	}
	return
}

func (h *HudLayer) Delete() {
	h.textRenderer.Delete()
}

func (h *HudLayer) Render() {
	// Render text for 'Geld', <Geld Amount>, 'Rating', <Rating Amount>
	hudItems := []string{strconv.Itoa(h.state.Rating), "YARPS", strconv.Itoa(h.state.Geld), "GELD"}

	var (
		textcache *twodee.TextCache
		texture   *twodee.Texture
		x         = h.camera.WorldBounds.Max.X()
		y         = h.camera.WorldBounds.Max.Y()
	)

	h.textRenderer.Bind()
	textcache = twodee.NewTextCache(h.regfont)

	for i, elem := range hudItems {
		textcache.SetText(elem)
		texture = textcache.Texture
		if texture != nil {
			if i%2 == 0 {
				x = x - (float32(texture.Width) + 10)
			} else {
				x = x - float32(texture.Width)
			}
			h.textRenderer.Draw(texture, x, y-float32(texture.Height))
		}
	}
	h.textRenderer.Unbind()
}

func (h *HudLayer) HandleEvent(evt twodee.Event) bool {
	return true
}

func (h *HudLayer) Reset() (err error) {
	if h.textRenderer != nil {
		h.textRenderer.Delete()
	}
	if h.textRenderer, err = twodee.NewTextRenderer(h.camera); err != nil {
		return
	}
	return
}

func (h *HudLayer) Update(elapsed time.Duration) {

}
