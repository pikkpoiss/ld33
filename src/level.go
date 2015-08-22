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
	"github.com/go-gl/mathgl/mgl32"
	"time"
)

type Level struct {
	Camera         *twodee.Camera
	Grid           *Grid
	Mobs           []*Mob
	ActiveMobCount int
	MousePos       mgl32.Vec2
	cursor         string
}

const (
	MaxMobs = 200
)

func NewLevel(sheet *twodee.Spritesheet) (level *Level, err error) {
	var (
		mobs   []*Mob
		grid   *Grid
		camera *twodee.Camera
	)
	mobs = make([]*Mob, MaxMobs)
	for i := 0; i < MaxMobs; i++ {
		mobs[i] = NewMob(sheet)
	}
	grid = NewGrid()
	if camera, err = twodee.NewCamera(
		twodee.Rect(0, 0, float32(grid.Width()), float32(grid.Height())),
		twodee.Rect(0, 0, 1024, 640),
	); err != nil {
		return
	}
	level = &Level{
		Camera:         camera,
		Grid:           grid,
		Mobs:           mobs,
		ActiveMobCount: 0,
	}
	return
}

func (l *Level) Update(elapsed time.Duration) {
	for i, mob := range l.Mobs {
		if !mob.Enabled { // No enabled mobs after first disabled mob.
			break
		}
		if mob.PendingDisable {
			l.disableMob(i)
		} else {
			mob.Update(elapsed, l)
		}
	}
}

func (l *Level) SetMouse(screenX, screenY float32) {
	x, y := l.Camera.ScreenToWorldCoords(screenX, screenY)
	l.MousePos = mgl32.Vec2{x, y}
}

func (l *Level) GetMouse() mgl32.Vec2 {
	return l.MousePos
}

func (l *Level) SetCursor(frame string) {
	l.cursor = frame
}

func (l *Level) GetCursor() string {
	return l.cursor
}

func (l *Level) SetBlock(pos mgl32.Vec2, block *Block) {
	var (
		gridCoords = l.Grid.WorldToGrid(pos)
	)
	if l.Grid.SetBlock(gridCoords, block) {
		l.Grid.CalculateDistances()
	}
}

func (l *Level) AddMob(pos mgl32.Vec2) {
	if l.ActiveMobCount == MaxMobs {
		// TODO: Do we need an error state?
		return
	}
	l.Mobs[l.ActiveMobCount].Activate(pos, 2.0)
	l.ActiveMobCount++
}

func (l *Level) disableMob(i int) {
	l.ActiveMobCount--
	l.Mobs[l.ActiveMobCount], l.Mobs[i] = l.Mobs[i], l.Mobs[l.ActiveMobCount]
	l.Mobs[l.ActiveMobCount].Disable()
}
