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
)

type UiState interface {
	HandleEvent(level *Level, evt twodee.Event) UiState
}

type NormalUiState struct {
}

func NewNormalUiState() UiState {
	return &NormalUiState{}
}

func (s *NormalUiState) HandleEvent(level *Level, evt twodee.Event) UiState {
	switch event := evt.(type) {
	case *twodee.MouseMoveEvent:
		level.SetMouse(event.X, event.Y)
	case *twodee.MouseButtonEvent:
		if event.Type == twodee.Press && event.Button == twodee.MouseButtonLeft {
			level.AddMob(level.GetMouse())
		}
	case *twodee.KeyEvent:
		if event.Type == twodee.Press && event.Code == twodee.Key1 {
			return NewBlockUiState()
		}
	}
	return nil
}

type BlockUiState struct {
}

func NewBlockUiState() UiState {
	return &BlockUiState{}
}

func (s *BlockUiState) HandleEvent(level *Level, evt twodee.Event) UiState {
	switch event := evt.(type) {
	case *twodee.MouseMoveEvent:
		level.SetMouse(event.X, event.Y)
	case *twodee.MouseButtonEvent:
		if event.Type == twodee.Press && event.Button == twodee.MouseButtonLeft {
			fmt.Printf("Drop block %v\n", level.GetMouse())
		}
	case *twodee.KeyEvent:
		if event.Type == twodee.Press && event.Code == twodee.Key0 {
			return NewNormalUiState()
		}
	}
	return nil
}
