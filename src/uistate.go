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
)

type UiState interface {
	Register(level *Level)
	Unregister(level *Level)
	HandleEvent(level *Level, evt twodee.Event) UiState
}

type BaseUiState struct {
}

func (s BaseUiState) HandleEvent(level *Level, evt twodee.Event) UiState {
	switch event := evt.(type) {
	case *twodee.MouseMoveEvent:
		level.SetMouse(event.X, event.Y)
	case *twodee.KeyEvent:
		if event.Type == twodee.Press {
			code := int(event.Code - twodee.Key1)
			if code >= 0 && code < 9 && code < len(HudBlocks) - 1 {
				return NewBlockUiState(HudBlocks[code + 1]) // Skip first (delete)
			}
			switch event.Code {
			case twodee.Key0:
				return NewNormalUiState()
			case twodee.KeyD:
				return NewDeleteUiState()
			}
		}
	}
	return nil
}

type NormalUiState struct {
	BaseUiState
}

func NewNormalUiState() UiState {
	return &NormalUiState{}
}

func (s *NormalUiState) Register(level *Level) {
	level.SetCursor("mouse_00")
}

func (s *NormalUiState) Unregister(level *Level) {
}

func (s *NormalUiState) HandleEvent(level *Level, evt twodee.Event) UiState {
	if state := s.BaseUiState.HandleEvent(level, evt); state != nil {
		return state
	}
	switch event := evt.(type) {
	case *twodee.MouseButtonEvent:
		if event.Type == twodee.Press && event.Button == twodee.MouseButtonLeft {
			if level.State.Debug {
				level.AddMob(level.GetMouse())
			}
		}
	}
	return nil
}

type DeleteUiState struct {
	BaseUiState
}

func NewDeleteUiState() UiState {
	return &DeleteUiState{}
}

func (s *DeleteUiState) Register(level *Level) {
	level.SetCursor("mouse_02")
}

func (s *DeleteUiState) Unregister(level *Level) {
	level.UnsetHighlights()
}

func (s *DeleteUiState) HandleEvent(level *Level, evt twodee.Event) UiState {
	if state := s.BaseUiState.HandleEvent(level, evt); state != nil {
		return state
	}
	switch event := evt.(type) {
	case *twodee.MouseMoveEvent:
		level.SetDeleteHighlights(level.GetMouse())
	case *twodee.MouseButtonEvent:
		if event.Type == twodee.Press && event.Button == twodee.MouseButtonLeft {
			level.DeleteBlock()
		}
	}
	return nil
}

type BlockUiState struct {
	BaseUiState
	target  *Block
	variant int
}

func NewBlockUiState(target *Block) UiState {
	return &BlockUiState{
		target:  target,
		variant: 0,
	}
}

func (s *BlockUiState) Register(level *Level) {
	level.SetCursor("mouse_01")
	level.SetHighlights(level.GetMouse(), s.target, s.variant)
}

func (s *BlockUiState) Unregister(level *Level) {
	level.UnsetHighlights()
}

func (s *BlockUiState) HandleEvent(level *Level, evt twodee.Event) UiState {
	if state := s.BaseUiState.HandleEvent(level, evt); state != nil {
		return state
	}
	switch event := evt.(type) {
	case *twodee.MouseMoveEvent:
		level.SetHighlights(level.GetMouse(), s.target, s.variant)
	case *twodee.MouseButtonEvent:
		if event.Type == twodee.Press && event.Button == twodee.MouseButtonLeft {
			if s.target.Cost <= level.State.Geld {
				level.State.Geld = level.State.Geld - s.target.Cost
				level.SetBlock(level.GetMouse(), s.target, s.variant)
				level.gameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayPlaceBlockEffect))
			}
		}
	case *twodee.KeyEvent:
		if event.Type == twodee.Press {
			switch event.Code {
			case twodee.KeyR:
				s.variant = (s.variant + 1) % len(s.target.Variants)
				level.SetHighlights(level.GetMouse(), s.target, s.variant)
			}
		}
	}
	return nil
}
