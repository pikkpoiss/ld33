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
)

func init() {
	// See https://code.google.com/p/go/issues/detail?id=3527
	runtime.LockOSThread()
}

type Application struct {
	Context *twodee.Context
}

func NewApplication() (app *Application, err error) {
	var (
		name      = "Twodee sample project"
		winbounds = twodee.Rect(0, 0, 1024, 640)
		context   *twodee.Context
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
	app = &Application{
		Context: context,
	}
	fmt.Printf("OpenGL version: %s\n", context.OpenGLVersion)
	fmt.Printf("Shader version: %s\n", context.ShaderVersion)
	return
}

func (a *Application) Draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func (a *Application) Delete() {
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
