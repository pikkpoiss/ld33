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
	"github.com/go-gl/mathgl/mgl32"
)

type State struct {
	Exit        bool
	Geld        int
	Rating      int
	Debug       bool
	MousePos    mgl32.Vec2
	MouseCursor string
	SplashState SplashState
}

func NewState() *State {
	return &State{
		Exit:        false,
		Geld:        100,
		Rating:      5,
		Debug:       false,
		MousePos:    mgl32.Vec2{0, 0},
		MouseCursor: "mouse_00",
		SplashState: SplashStart,
	}
}
