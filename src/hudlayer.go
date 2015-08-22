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
	"time"
)

type HudLayer struct {
	camera       *twodee.Camera
	cameraBounds twodee.Rectangle
	app          *Application
}

func NewHudLayer(app *Application) (layer *HudLayer, err error) {
	var (
		camera *twodee.Camera
	)
	if camera, err = twodee.NewCamera(
		twodee.Rect(0, 0, float32(GridWidth), float32(GridHeight)),
		twodee.Rect(0, 0, ScreenWidth, ScreenHeight),
	); err != nil {
		return
	}
	layer = &HudLayer{
		camera: camera,
		app:    app,
	}
	return
}

func (h *HudLayer) Delete() {
}

func (h *HudLayer) Render() {

}

func (h *HudLayer) HandleEvent(evt twodee.Event) bool {
	return true
}

func (h *HudLayer) Reset() (err error) {
	h.Delete()
	return
}

func (h *HudLayer) Update(elapsed time.Duration) {

}
