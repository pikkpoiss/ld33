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
	"time"
)

type GridItem struct {
	passable  bool
	distance  int32
	frame     string
	animation *twodee.FrameAnimation
	frames    BlockAnimations
	state     BlockState
}

func NewGridItem(passable bool, frame string, frames BlockAnimations) *GridItem {
	var (
		animation *twodee.FrameAnimation
		state     = BlockNormal
	)
	if frames != nil {
		animation = twodee.NewFrameAnimation(
			100*time.Millisecond,
			frames[state],
		)
	}
	return &GridItem{
		passable:  passable,
		distance:  -1,
		frame:     frame,
		frames:    frames,
		animation: animation,
		state:     state,
	}
}

func (i *GridItem) SetState(state BlockState) {
	if state != i.state {
		if frames, ok := i.frames[state]; ok {
			i.animation.SetSequence(frames)
			i.state = state
		}
	}
}

func (i *GridItem) Update(elapsed time.Duration) {
	if i.animation != nil {
		i.animation.Update(elapsed)
	}
}

func (i *GridItem) Frame() string {
	if i.animation != nil {
		return fmt.Sprintf(i.frame, i.animation.Current)
	}
	return i.frame
}

func (i *GridItem) Passable() bool {
	return i.passable
}

func (i *GridItem) Opaque() bool {
	return i.passable
}

func (i *GridItem) SetDistance(dist int32) {
	i.distance = dist
}

func (i *GridItem) Distance() int32 {
	return i.distance
}
