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
	"github.com/go-gl/gl/v3.3-core/gl"
	"runtime"
	"time"
)

func init() {
	// See https://code.google.com/p/go/issues/detail?id=3527
	runtime.LockOSThread()
}

type Application struct {
	layers           *twodee.Layers
	Context          *twodee.Context
	State            *State
	GameEventHandler *twodee.GameEventHandler
	AudioSystem      *AudioSystem
	gameLayer        *GameLayer
}

func NewApplication() (app *Application, err error) {
	var (
		name             = "Screamporium - LD33"
		layers           *twodee.Layers
		context          *twodee.Context
		menulayer        *MenuLayer
		hudlayer         *HudLayer
		splashlayer      *SplashLayer
		winbounds        = twodee.Rect(0, 0, ScreenWidth, ScreenHeight)
		state            = NewState()
		gameEventHandler = twodee.NewGameEventHandler(NumGameEventTypes)
		audioSystem      *AudioSystem
	)
	if context, err = twodee.NewContext(); err != nil {
		return
	}
	context.SetFullscreen(false)
	context.SetCursor(false)
	if err = context.CreateWindow(
		int(winbounds.Max.X()),
		int(winbounds.Max.Y()),
		name,
	); err != nil {
		return
	}
	context.SetSwapInterval(2)
	layers = twodee.NewLayers()
	app = &Application{
		layers:           layers,
		Context:          context,
		State:            state,
		GameEventHandler: gameEventHandler,
	}
	if app.gameLayer, err = NewGameLayer(state, app); err != nil {
		return
	}
	layers.Push(app.gameLayer)
	if splashlayer, err = NewSplashLayer(state, app, app.gameLayer.level.Grid); err != nil {
		return
	}
	layers.Push(splashlayer)
	if menulayer, err = NewMenuLayer(winbounds, state, app); err != nil {
		return
	}
	layers.Push(menulayer)
	if hudlayer, err = NewHudLayer(state, app.gameLayer.level.Grid, app); err != nil {
		return
	}
	layers.Push(hudlayer)
	if audioSystem, err = NewAudioSystem(app); err != nil {
		return
	}
	app.AudioSystem = audioSystem
	fmt.Printf("OpenGL version: %s\n", context.OpenGLVersion)
	fmt.Printf("Shader version: %s\n", context.ShaderVersion)
	return
}

func (a *Application) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	a.layers.Render()
}

func (a *Application) Update(elapsed time.Duration) {
	a.layers.Update(elapsed)
}

func (a *Application) Delete() {
	a.layers.Delete()
	a.Context.Delete()
	a.AudioSystem.Delete()
}

func (a *Application) SetUiState(state UiState) {
	a.gameLayer.SetUiState(state)
}

func (a *Application) UnsetHighlights() {
	a.gameLayer.UnsetHighlights()
}

func (a *Application) ProcessEvents() {
	var (
		evt   twodee.Event
		loop  = true
		count = 0
	)
	for loop {
		select {
		case evt = <-a.Context.Events.Events:
			a.layers.HandleEvent(evt)
			count++
			if count > 10 {
				loop = false
			}
		default:
			loop = false
		}
	}
}

func main() {
	var (
		app *Application
		err error
	)

	if app, err = NewApplication(); err != nil {
		panic(err)
	}
	defer app.Delete()

	var (
		current_time = time.Now()
		updated_to   = current_time
		step         = twodee.Step30Hz
	)
	for !app.Context.ShouldClose() && !app.State.Exit {
		app.Context.Events.Poll()
		app.GameEventHandler.Poll()
		app.ProcessEvents()
		for !updated_to.After(current_time) {
			app.Update(step)
			updated_to = updated_to.Add(step)
		}
		current_time = current_time.Add(step)
		time.Sleep(current_time.Sub(time.Now()))
		app.Draw()
		app.Context.SwapBuffers()
	}
}
