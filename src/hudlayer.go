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
	"image/color"
	"io/ioutil"
	"strconv"
	"time"
)

type HudItem struct {
	HitBox      twodee.Rectangle
	Enabled     bool
	Highlighted bool
	Block       *Block
}

type HudLayer struct {
	camera         *twodee.Camera
	textRenderer   *twodee.TextRenderer
	regFont        *twodee.FontFace
	pixelFont      *twodee.FontFace
	spriteSheet    *twodee.Spritesheet
	spriteTexture  *twodee.Texture
	spriteRenderer *twodee.SpriteRenderer
	state          *State
	app            *Application
	textCache      map[string]*twodee.TextCache
	items          []HudItem
	textScale      float32
}

func NewHudLayer(state *State, grid *Grid, app *Application) (layer *HudLayer, err error) {
	var (
		regFont        *twodee.FontFace
		pixelFont      *twodee.FontFace
		bg             = color.Transparent
		regFontPath    = "resources/fonts/Prototype.ttf"
		pixelFontPath  = "resources/fonts/slkscr.ttf"
		camera         *twodee.Camera
		regFontColor           = color.RGBA{0, 0, 0, 255}
		pixelFontColor         = color.RGBA{0, 0, 0, 255}
		textScale      float32 = 1.0 / 32.0
		textSize               = 32.0
	)
	if camera, err = twodee.NewCamera(
		twodee.Rect(0, 0, float32(grid.Width()), float32(grid.Height())),
		//twodee.Rect(0, 0, ScreenWidth, ScreenHeight),
		twodee.Rect(0, 0, ScreenWidth, ScreenHeight),
	); err != nil {
		return
	}
	if regFont, err = twodee.NewFontFace(regFontPath, textSize, regFontColor, bg); err != nil {
		return
	}
	if pixelFont, err = twodee.NewFontFace(pixelFontPath, textSize, pixelFontColor, bg); err != nil {
		return
	}
	layer = &HudLayer{
		camera:    camera,
		regFont:   regFont,
		pixelFont: pixelFont,
		state:     state,
		app:       app,
		textCache: map[string]*twodee.TextCache{},
		textScale: textScale,
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
		configs   = []twodee.SpriteConfig{}
		texture   *twodee.Texture
		xText     = h.camera.WorldBounds.Max.X()
		yText     = h.camera.WorldBounds.Max.Y()
		texHeight float32
		texWidth  float32
		//ySprite         = h.spriteCamera.WorldBounds.Max.Y()
		//verticalSpacing = 80
	)

	// Render toolbar for selecting blocks to place
	h.spriteTexture.Bind()

	for _, item := range h.items {
		x := item.HitBox.Max.X()
		y := item.HitBox.Min.Y()
		switch {
		case item.Highlighted:
			configs = append(configs, h.toolbarSpriteConfig(h.spriteSheet, 4, y))
			configs = append(configs, h.highlightSpriteConfig(h.spriteSheet, mgl32.Vec2{x, y}, "highlight_00"))
		case item.Block.Cost <= h.state.Geld:
			configs = append(configs, h.toolbarSpriteConfig(h.spriteSheet, 1, y))
		default:
			configs = append(configs, h.toolbarSpriteConfig(h.spriteSheet, 15, y))
		}
	}

	configs = append(
		configs,
		h.cursorSpriteConfig(h.spriteSheet, h.state.MousePos, h.state.MouseCursor),
	)

	if len(configs) > 0 {
		h.spriteRenderer.Draw(configs)
	}
	h.spriteTexture.Unbind()

	// Put text on top

	h.textRenderer.Bind()

	// Render text for toolbar
	texture = h.cacheText("toolbar", h.regFont, "Toolbar")
	if texture != nil {
		h.textRenderer.Draw(texture, 5, yText-float32(texture.Height), h.textScale)
	}

	// Render text for 'Geld', <Geld Amount>, 'Rating', <Rating Amount>
	for i, elem := range hudItems {
		texture = h.cacheText(fmt.Sprintf("toolbar%v", i), h.regFont, elem)
		if texture != nil {
			texHeight = float32(texture.Height) * h.textScale
			texWidth = float32(texture.Width) * h.textScale
			if i%2 == 0 {
				xText = xText - (texWidth + 1)
			} else {
				xText = xText - texWidth
			}
			h.textRenderer.Draw(texture, xText, yText-texHeight, h.textScale)
		}
	}

	for i, item := range h.items {
		texture = h.cacheText(fmt.Sprintf("key%v", i), h.pixelFont, strconv.Itoa(i+1))
		if texture != nil {
			h.textRenderer.Draw(texture, item.HitBox.Min.X(), item.HitBox.Min.Y(), h.textScale)
		}
		if item.Highlighted {
			texture = h.cacheText("highlight", h.pixelFont, item.Block.Title)
			h.textRenderer.Draw(texture, item.HitBox.Max.X()+1, item.HitBox.Min.Y(), h.textScale)
		}
	}

	h.textRenderer.Unbind()
}

func (h *HudLayer) HandleEvent(evt twodee.Event) bool {
	switch event := evt.(type) {
	case *twodee.MouseButtonEvent:
		if event.Type == twodee.Press && event.Button == twodee.MouseButtonLeft {
			for _, item := range h.items {
				if item.Highlighted {
					h.app.SetUiState(NewBlockUiState(item.Block))
					return false
				}
			}
		}
	}
	return true
}

func (h *HudLayer) makeItems() {
	var (
		yMax      = h.camera.WorldBounds.Max.Y()
		blocks    = []*Block{&SkellyBlock, &SpikesBlock, &CornerBlock}
		block     *Block
		boxHeight float32 = 2
		boxWidth  float32 = 2
		boxOffset float32 = 4
		bottom    float32
		top       float32
		i         int
	)
	h.items = make([]HudItem, len(blocks))
	for i, block = range blocks {
		top = yMax - (boxHeight*float32(i) + boxOffset)
		bottom = top - boxHeight
		h.items[i].Enabled = false
		h.items[i].Highlighted = false
		h.items[i].HitBox = twodee.Rect(0, bottom, boxWidth, top)
		h.items[i].Block = block
	}
}

func (h *HudLayer) Reset() (err error) {
	if h.textRenderer != nil {
		h.textRenderer.Delete()
	}
	if h.textRenderer, err = twodee.NewTextRenderer(h.camera); err != nil {
		return
	}
	if h.spriteRenderer, err = twodee.NewSpriteRenderer(h.camera); err != nil {
		return
	}
	if err = h.loadSpritesheet(); err != nil {
		return
	}
	h.makeItems()
	return
}

func (h *HudLayer) Update(elapsed time.Duration) {
	var overlaps bool
	for i, item := range h.items {
		overlaps = item.HitBox.ContainsPoint(twodee.Point{h.state.MousePos})
		h.items[i].Highlighted = overlaps
		if overlaps {
			h.app.UnsetHighlights()
		}
	}
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
	var frame *twodee.SpritesheetFrame
	frame = sheet.GetFrame(fmt.Sprintf("numbered_squares_%02v", block))
	xPosition := (frame.Width / 2.0) + 1.2
	yPosition := y + (frame.Height / 2.0) // Bottom aligned
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			xPosition, yPosition, 0,
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}
}

func (h *HudLayer) cursorSpriteConfig(sheet *twodee.Spritesheet, pt mgl32.Vec2, cursor string) twodee.SpriteConfig {
	frame := sheet.GetFrame(cursor)
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			pt.X(), pt.Y(), 0.0,
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}
}

func (h *HudLayer) highlightSpriteConfig(sheet *twodee.Spritesheet, pt mgl32.Vec2, name string) twodee.SpriteConfig {
	frame := sheet.GetFrame(name)
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			pt.X() + frame.Width/2.0, pt.Y() + frame.Height/6.0, 0.0, // Left aligned
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}
}
