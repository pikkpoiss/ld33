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
	"github.com/go-gl/mathgl/mgl32"
	"time"
)

type MobState int32

const (
	_                = iota
	Walking MobState = 1 << iota
	Left
	Right
)

var MobAnimations = map[MobState][]int{
	Walking | Right: []int{0, 1, 2, 3, 4, 5, 6, 7},
	Walking | Left:  []int{0, 1, 2, 3, 4, 5, 6, 7},
}

type Mob struct {
	*twodee.AnimatingEntity
	State          MobState
	Speed          float32
	Fear           float64
	Enabled        bool
	PendingDisable bool
	Pos            mgl32.Vec2
}

func NewMob(sheet *twodee.Spritesheet) *Mob {
	var (
		frame = sheet.GetFrame("human01_00")
	)
	return &Mob{
		AnimatingEntity: twodee.NewAnimatingEntity(
			0, 0,
			frame.Width, frame.Height,
			0.0,
			twodee.Step10Hz,
			MobAnimations[Walking|Right],
		),
		Fear: 1.0,
	}
}

func (m *Mob) Update(elapsed time.Duration, level *Level) {
	m.AnimatingEntity.Update(elapsed)
	m.moveTowardExit(elapsed, level)
}

func (m *Mob) moveTowardExit(elapsed time.Duration, level *Level) {
	var (
		dest     mgl32.Vec2
		ok       bool
		pct      = float32(elapsed) / float32(time.Second)
		gridDist mgl32.Vec2
		goalDist int32
		stepDist = pct * m.Speed
	)
	if dest, goalDist, ok = level.Grid.GetNextStepToSink(m.Pos); !ok {
		return
	}
	gridDist = dest.Sub(m.Pos)
	if goalDist == 1 && gridDist.Len() < stepDist+0.5 {
		m.PendingDisable = true
	}
	if gridDist.X() > 0 {
		m.swapState(Left, Right)
	} else {
		m.swapState(Right, Left)
	}
	m.Pos = m.Pos.Add(gridDist.Normalize().Mul(stepDist))
}

func (m *Mob) Activate(pos mgl32.Vec2, speed float32) {
	m.Enabled = true
	m.PendingDisable = false
	m.Pos = pos
	m.Speed = speed
	m.State = Walking | Right
}

func (m *Mob) Disable() {
	m.Enabled = false
	m.PendingDisable = false
	m.Fear = 1.0
}

func (m *Mob) AddSpriteConfig(sheet *twodee.Spritesheet, config []twodee.SpriteConfig) []twodee.SpriteConfig {
	var (
		frame               = sheet.GetFrame(fmt.Sprintf("human01_%02d", m.Frame()))
		scaleX      float32 = 1.0
		view        twodee.ModelViewConfig
		overlayview twodee.ModelViewConfig
	)
	if m.State&Left == Left {
		scaleX = -1.0
	}
	view = twodee.ModelViewConfig{
		m.Pos.X(), m.Pos.Y() + frame.Height/4.0, 0.0,
		0, 0, 0,
		scaleX, 1.0, 1.0,
	}
	overlayview = twodee.ModelViewConfig{
		m.Pos.X(), m.Pos.Y() + frame.Height/4.0 - 0.01, 0.0,
		0, 0, 0,
		scaleX, 1.0, 1.0,
	}
	config = append(config, twodee.SpriteConfig{
		View:  view,
		Frame: frame.Frame,
	})
	switch {
	case m.Fear < 5:
		config = append(config, twodee.SpriteConfig{
			View:  overlayview,
			Frame: sheet.GetFrame("overlays_01").Frame,
		})
	case m.Fear > 9:
		config = append(config, twodee.SpriteConfig{
			View:  overlayview,
			Frame: sheet.GetFrame("overlays_02").Frame,
		})
	case m.Fear > 7:
		config = append(config, twodee.SpriteConfig{
			View:  overlayview,
			Frame: sheet.GetFrame("overlays_00").Frame,
		})
	}
	return config
}

func (m *Mob) remState(state MobState) {
	m.setState(m.State & ^state)
}

func (m *Mob) addState(state MobState) {
	m.setState(m.State | state)
}

func (m *Mob) swapState(rem, add MobState) {
	m.setState(m.State & ^rem | add)
}

func (m *Mob) setState(state MobState) {
	if state != m.State {
		m.State = state
		if frames, ok := MobAnimations[m.State]; ok {
			m.SetFrames(frames)
		}
	}
}

// IncreaseFear increments the mob's fear counter and returns a bool indicating
// whether the mob is still alive or has passed away from fright.
func (m *Mob) IncreaseFear(fear float64) bool {
	m.Fear += fear
	return m.Fear < 10
}
