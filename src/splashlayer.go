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

type SplashState int32

const (
	_                          = iota
	SplashDisabled SplashState = 1 << iota
	SplashStart
	SplashWin
	SplashLose
)

type SplashLayer struct {
	app            *Application
	state          *State
	splashRenderer *SplashRenderer
	camera         *twodee.Camera
}

func NewSplashLayer(state *State, app *Application, grid *Grid) (layer *SplashLayer, err error) {
	var (
		camera *twodee.Camera
	)
	if camera, err = twodee.NewCamera(
		twodee.Rect(0, 0, float32(grid.Width()), float32(grid.Height())),
		twodee.Rect(0, 0, ScreenWidth, ScreenHeight),
	); err != nil {
		return
	}
	layer = &SplashLayer{
		app:    app,
		state:  state,
		camera: camera,
	}
	err = layer.Reset()
	return
}

func (l *SplashLayer) Delete() {
	l.splashRenderer.Delete()
}

func (l *SplashLayer) Update(elapsed time.Duration) {
}

func (l *SplashLayer) Render() {
	if l.state.SplashState != SplashDisabled {
		l.splashRenderer.Draw(l.state)
	}
}

func (l *SplashLayer) HandleEvent(evt twodee.Event) bool {
	if l.state.SplashState != SplashDisabled {
		switch event := evt.(type) {
		case *twodee.KeyEvent:
			if event.Type != twodee.Press {
				break
			}
			switch event.Code {
			case twodee.KeySpace:
				l.state.SplashState = SplashDisabled
			}
		case *twodee.MouseButtonEvent:
			if event.Type == twodee.Press && event.Button == twodee.MouseButtonLeft {
				l.state.SplashState = SplashDisabled
			}
		}
	}
	return true
}

func (l *SplashLayer) Reset() (err error) {
	l.splashRenderer, err = NewSplashRenderer(l.camera)
	return
}
