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
	"math"
	"time"
)

type SpawnZone struct {
	Pos    Ivec2
	charge float64
}

func (s *SpawnZone) AddCharge(c float64) {
	s.charge += c
}

func NewSpawnZone(p Ivec2) SpawnZone {
	s := SpawnZone{}
	s.Pos = p
	return s
}

// Spawn checks if the SpawnZone has accumulated enough charge to spawn a unit.
// If so, it returns true and removes the requisite amount of charge from the
// zone. Otherwise, returns false.
func (s *SpawnZone) Spawn() bool {
	remainingCharge := s.charge - 1
	if remainingCharge > 0 {
		s.charge = remainingCharge
	}
	return remainingCharge > 0
}

type Highlight struct {
	Pos   Ivec2
	Frame string
}

type Level struct {
	Camera         *twodee.Camera
	Grid           *Grid
	State          *State
	Mobs           []Mob
	ActiveMobCount int
	MousePos       mgl32.Vec2
	Highlights     []Highlight
	cursor         string
	entries        []SpawnZone
	exit           SpawnZone
	blocks         map[Ivec2]*Block
	fearHistory    []float64
	fearIndex      int
}

const (
	MaxMobs = 200
)

func NewLevel(state *State, sheet *twodee.Spritesheet) (level *Level, err error) {
	var (
		mobs    = make([]Mob, MaxMobs)
		grid    *Grid
		camera  *twodee.Camera
		entries = []SpawnZone{
			NewSpawnZone(Ivec2{4, 9}),
			NewSpawnZone(Ivec2{4, 14}),
			NewSpawnZone(Ivec2{4, 4}),
		}
		exit        = NewSpawnZone(Ivec2{20, 10})
		fearHistory = make([]float64, 100)
	)
	if grid, err = NewGrid(); err != nil {
		return
	}
	for _, entry := range entries {
		grid.AddSource(entry.Pos)
	}
	grid.SetSink(exit.Pos)
	grid.CalculateDistances()

	for i := 0; i < MaxMobs; i++ {
		mobs[i] = *NewMob(sheet)
	}
	if camera, err = twodee.NewCamera(
		twodee.Rect(0, 0, float32(grid.Width()), float32(grid.Height())),
		twodee.Rect(0, 0, ScreenWidth, ScreenHeight),
	); err != nil {
		return
	}

	for i := 0; i < 100; i++ {
		fearHistory[i] = 5
	}

	level = &Level{
		Camera:         camera,
		Grid:           grid,
		State:          state,
		Mobs:           mobs,
		ActiveMobCount: 0,
		entries:        entries,
		exit:           exit,
		blocks:         make(map[Ivec2]*Block),
		fearHistory:    fearHistory,
		fearIndex:      0,
	}
	return
}

func (l *Level) updateMobs(elapsed time.Duration) {
	for i := range l.Mobs {
		mob := &l.Mobs[i]
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

func (l *Level) updateSpawns(elapsed time.Duration) {
	// TODO: Calculate amount of charge as f(elapsed, rating)
	charge := 0.002 * float64(l.State.Rating)
	for i := range l.entries {
		entry := &l.entries[i]
		entry.AddCharge(charge)
		for entry.Spawn() {
			l.SpawnMob(entry.Pos)
		}
	}
}

func (l *Level) updateBlocks(elapsed time.Duration) {
	for pos, block := range l.blocks {
		posV := mgl32.Vec2{float32(pos.X()), float32(pos.Y())}
		fear := block.FearPerSec * elapsed.Seconds()
		numHit := 0
		for i := range l.Mobs {
			if numHit >= block.MaxTargets {
				break
			}
			mob := &l.Mobs[i]
			if mob.Pos.Sub(posV).Len() <= block.Range {
				numHit++
				mob.IncreaseFear(fear)
			}
		}
	}
}

// Update computes a new simulation step for this level.
func (l *Level) Update(elapsed time.Duration) {
	l.updateBlocks(elapsed)
	l.updateMobs(elapsed)
	l.updateSpawns(elapsed)
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
	if center, ok := l.Grid.SetBlock(gridCoords, block); ok {
		l.blocks[center] = block
		l.Grid.CalculateDistances()
	}
}

// calculateRating returns the rounded integer average of all values in
// fearHistory.
// TODO: The rating should probably be cached and recalculated when mobs
// despawn, as opposed to generated anew each time like this...
func (l *Level) calculateRating() int {
	sum := 0.0
	for _, v := range l.fearHistory {
		sum += v
	}
	return int(math.Floor(sum/float64(len(l.fearHistory)) + 0.5))
}

func (l *Level) ClearHighlights() {
	l.Highlights = l.Highlights[0:0]
}

func (l *Level) SetHighlights(pos mgl32.Vec2, block *Block) {
	var (
		pre   = l.Grid.WorldToGrid(pos)
		post  = pre.Plus(block.Offset)
		frame = "special_squares_02"
	)
	l.ClearHighlights()
	if !l.Grid.IsBlockValid(pre, block) {
		frame = "special_squares_03"
	}
	for y := 0; y < len(block.Template); y++ {
		for x := 0; x < len(block.Template[y]); x++ {
			l.Highlights = append(l.Highlights, Highlight{
				post.Plus(Ivec2{int32(x), int32(y)}),
				frame,
			})
		}
	}
}

func (l *Level) SpawnMob(v Ivec2) {
	p := mgl32.Vec2{float32(v.X()), float32(v.Y())}
	l.AddMob(p)
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
	l.fearHistory[l.fearIndex] = l.Mobs[i].Fear
	l.fearIndex = (l.fearIndex + 1) % len(l.fearHistory)
	l.State.Rating = l.calculateRating()
	l.State.Geld += int(math.Floor(l.Mobs[i].Fear*10.0 + 0.5))
	l.ActiveMobCount--
	l.Mobs[l.ActiveMobCount], l.Mobs[i] = l.Mobs[i], l.Mobs[l.ActiveMobCount]
	l.Mobs[l.ActiveMobCount].Disable()
}
