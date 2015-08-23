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
	"image/color"
	"io/ioutil"
	"strconv"
	"time"
)

type HudLayer struct {
	textCamera     *twodee.Camera
	spriteCamera   *twodee.Camera
	textRenderer   *twodee.TextRenderer
	regFont        *twodee.FontFace
	pixelFont      *twodee.FontFace
	spriteSheet    *twodee.Spritesheet
	spriteTexture  *twodee.Texture
	spriteRenderer *twodee.SpriteRenderer
	state          *State
	app            *Application
	textCache      map[string]*twodee.TextCache
}

func NewHudLayer(state *State, grid *Grid, app *Application) (layer *HudLayer, err error) {
	var (
		regFont        *twodee.FontFace
		pixelFont      *twodee.FontFace
		bg             = color.Transparent
		regFontPath    = "resources/fonts/Prototype.ttf"
		pixelFontPath  = "resources/fonts/slkscr.ttf"
		textCamera     *twodee.Camera
		spriteCamera   *twodee.Camera
		regFontColor   = color.RGBA{0, 0, 0, 255}
		pixelFontColor = color.RGBA{0, 0, 0, 255}
	)
	if textCamera, err = twodee.NewCamera(
		twodee.Rect(0, 0, ScreenWidth, ScreenHeight),
		twodee.Rect(0, 0, ScreenWidth, ScreenHeight),
	); err != nil {
		return
	}
	if spriteCamera, err = twodee.NewCamera(
		twodee.Rect(0, 0, float32(grid.Width()), float32(grid.Height())),
		twodee.Rect(0, 0, ScreenWidth, ScreenHeight),
	); err != nil {
		return
	}
	if regFont, err = twodee.NewFontFace(regFontPath, 32, regFontColor, bg); err != nil {
		return
	}
	if pixelFont, err = twodee.NewFontFace(pixelFontPath, 32, pixelFontColor, bg); err != nil {
		return
	}
	layer = &HudLayer{
		textCamera:   textCamera,
		spriteCamera: spriteCamera,
		regFont:      regFont,
		pixelFont:    pixelFont,
		state:        state,
		app:          app,
		textCache:    map[string]*twodee.TextCache{},
	}
	err = layer.Reset()
	return
}

func (h *HudLayer) Delete() {
	h.textRenderer.Delete()
	h.spriteRenderer.Delete()
}

func (h *HudLayer) cacheText(key string, font *twodee.FontFace, value string) *twodee.Texture {
	var (
		ok    bool
		cache *twodee.TextCache
	)
	if cache, ok = h.textCache[key]; !ok {
		cache = twodee.NewTextCache(font)
		h.textCache[key] = cache
	}
	cache.SetText(value)
	return cache.Texture
}

func (h *HudLayer) Render() {
	hudItems := []string{strconv.Itoa(h.state.Rating), "RATING", strconv.Itoa(h.state.Geld), "GELD"}

	var (
		configs         = []twodee.SpriteConfig{}
		texture         *twodee.Texture
		xText           = h.textCamera.WorldBounds.Max.X()
		yText           = h.textCamera.WorldBounds.Max.Y()
		ySprite         = h.spriteCamera.WorldBounds.Max.Y()
		verticalSpacing = 80
	)

	h.textRenderer.Bind()

	// Render text for toolbar
	texture = h.cacheText("toolbar", h.regFont, "Toolbar")
	if texture != nil {
		h.textRenderer.Draw(texture, 5, yText-float32(texture.Height))
	}

	// Render text for 'Geld', <Geld Amount>, 'Rating', <Rating Amount>
	for i, elem := range hudItems {
		texture = h.cacheText(fmt.Sprintf("toolbar%v", i), h.regFont, elem)
		if texture != nil {
			if i%2 == 0 {
				xText = xText - (float32(texture.Width) + 10)
			} else {
				xText = xText - float32(texture.Width)
			}
			h.textRenderer.Draw(texture, xText, yText-float32(texture.Height))
		}
	}

	// Render text for each of the available blocks to purchase
	blocks := []*Block{&SkellyBlock, &SpikesBlock, &CornerBlock}
	blockCost := 0
	for i, block := range blocks {
		blockCost = block.Cost
		if blockCost <= h.state.Geld {
			texture = h.cacheText(fmt.Sprintf("key%v", i), h.pixelFont, strconv.Itoa(i+1))
			if texture != nil {
				h.textRenderer.Draw(texture, 5, yText-float32(verticalSpacing*(i+1)))
			}
		}
	}

	h.textRenderer.Unbind()

	// Render toolbar for selecting blocks to place
	h.spriteTexture.Bind()
	for blockNum, block := range blocks {
		blockCost = block.Cost
		if blockCost <= h.state.Geld {
			configs = append(configs, h.toolbarSpriteConfig(h.spriteSheet, float32(blockNum), ySprite))
		} else {
			configs = append(configs, h.toolbarSpriteConfig(h.spriteSheet, 15, ySprite))
		}
	}
	h.spriteRenderer.Draw(configs)
	h.spriteTexture.Unbind()

}

func (h *HudLayer) HandleEvent(evt twodee.Event) bool {
	return true
}

func (h *HudLayer) Reset() (err error) {
	if h.textRenderer != nil {
		h.textRenderer.Delete()
	}
	if h.textRenderer, err = twodee.NewTextRenderer(h.textCamera); err != nil {
		return
	}
	if h.spriteRenderer, err = twodee.NewSpriteRenderer(h.spriteCamera); err != nil {
		return
	}
	if err = h.loadSpritesheet(); err != nil {
		return
	}
	return
}

func (h *HudLayer) Update(elapsed time.Duration) {

}

func (h *HudLayer) loadSpritesheet() (err error) {
	var (
		data []byte
	)
	if data, err = ioutil.ReadFile("resources/spritesheet.json"); err != nil {
		return
	}
	if h.spriteSheet, err = twodee.ParseTexturePackerJSONArrayString(
		string(data),
		PxPerUnit,
	); err != nil {
		return
	}
	if h.spriteTexture, err = twodee.LoadTexture(
		"resources/"+h.spriteSheet.TexturePath,
		twodee.Nearest,
	); err != nil {
		return
	}
	return
}

func (h *HudLayer) toolbarSpriteConfig(sheet *twodee.Spritesheet, block float32, y float32) twodee.SpriteConfig {
	var toolbarSpriteSpacing float32 = 1.4
	var toolbarSpriteVerticalOffset float32 = 2.2
	var frame *twodee.SpritesheetFrame
	frame = sheet.GetFrame(fmt.Sprintf("numbered_squares_%02v", block))
	xPosition := (frame.Width / 2.0) + 1.2
	yPosition := y - (block * (frame.Height + toolbarSpriteSpacing)) - toolbarSpriteVerticalOffset
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			xPosition, yPosition, 0,
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}
}
