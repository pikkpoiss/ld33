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
	"time"
)

type Mob struct {
	Pos            mgl32.Vec2
	Speed          float32
	Enabled        bool
	PendingDisable bool
}

func (m *Mob) Update(elapsed time.Duration, level *Level) {
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
	if dest, goalDist, ok = level.Grid.GetNextStepToExit(m.Pos); !ok {
		return
	}
	gridDist = dest.Sub(m.Pos)
	if goalDist == 1 && gridDist.Len() < stepDist + 0.5 {
		m.PendingDisable = true
	}
	m.Pos = m.Pos.Add(gridDist.Normalize().Mul(stepDist))
}

func (m *Mob) Activate(pos mgl32.Vec2, speed float32) {
	m.Enabled = true
	m.PendingDisable = false
	m.Pos = pos
	m.Speed = speed
}

func (m *Mob) Disable() {
	m.Enabled = false
	m.PendingDisable = false
}
