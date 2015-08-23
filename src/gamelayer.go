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

type GameLayer struct {
	gameRenderer  *GameRenderer
	spriteSheet   *twodee.Spritesheet
	spriteTexture *twodee.Texture
	app           *Application
	level         *Level
	uiState       UiState
	state         *State
}

func NewGameLayer(state *State, app *Application) (layer *GameLayer, err error) {
	layer = &GameLayer{
		app:   app,
		state: state,
	}
	err = layer.Reset()
	return
}

func (l *GameLayer) Delete() {
	l.gameRenderer.Delete()
}

func (l *GameLayer) Render() {
	l.spriteTexture.Bind()
	l.gameRenderer.Draw(l.level)
	l.spriteTexture.Unbind()
}

func (l *GameLayer) SetUiState(state UiState) {
	l.uiState.Unregister(l.level)
	l.uiState = state
	l.uiState.Register(l.level)
}

func (l *GameLayer) HandleEvent(evt twodee.Event) bool {
	if newState := l.uiState.HandleEvent(l.level, evt); newState != nil {
		l.SetUiState(newState)
	}

	switch event := evt.(type) {
	case *twodee.KeyEvent:
		if event.Type == twodee.Release {
			break
		}
		switch event.Code {
		case twodee.KeyM:
			if twodee.MusicIsPaused() {
				l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(ResumeMusic))
			} else {
				l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PauseMusic))
			}
		}
	}
	return true
}

func (l *GameLayer) Reset() (err error) {
	if err = l.loadSpritesheet(); err != nil {
		return
	}
	if l.level, err = NewLevel(l.state, l.spriteSheet); err != nil {
		return
	}
	l.uiState = NewNormalUiState()
	l.uiState.Register(l.level)
	if l.gameRenderer, err = NewGameRenderer(l.level, l.spriteSheet); err != nil {
		return
	}
	l.app.GameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayBackgroundMusic))
	return
}

func (l *GameLayer) Update(elapsed time.Duration) {
	l.level.Update(elapsed)
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
