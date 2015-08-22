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
	layers  *twodee.Layers
	Context *twodee.Context
	State   *State
}

func NewApplication() (app *Application, err error) {
	var (
		name      = "Twodee sample project"
		layers    *twodee.Layers
		context   *twodee.Context
		menulayer *MenuLayer
		gamelayer *GameLayer
		winbounds = twodee.Rect(0, 0, 1024, 640)
		state     = NewState()
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
	layers = twodee.NewLayers()
	app = &Application{
		layers:  layers,
		Context: context,
		State: state,
	}
	if gamelayer, err = NewGameLayer(); err != nil {
		return
	}
	layers.Push(gamelayer)
	if menulayer, err = NewMenuLayer(winbounds, state, app); err != nil {
		return
	}
	layers.Push(menulayer)
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

	for !app.Context.ShouldClose() {
		app.Context.Events.Poll()
		app.Draw()
		app.Context.SwapBuffers()
	}
}
