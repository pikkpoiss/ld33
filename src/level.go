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
	"sort"
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

type BlockPlacement struct {
	Block   *Block
	Variant int
}

type Level struct {
	Camera           *twodee.Camera
	Grid             *Grid
	State            *State
	Mobs             []Mob
	Decals           []*Decal
	ActiveMobCount   int
	ActiveDecalCount int
	Highlights       []Highlight
	entries          []SpawnZone
	exit             SpawnZone
	blocks           map[Ivec2]BlockPlacement
	fearHistory      []float64
	fearIndex        int
}

const (
	MaxMobs   = 200
	MaxDecals = 10
)

func NewLevel(state *State, sheet *twodee.Spritesheet) (level *Level, err error) {
	var (
		mobs    = make([]Mob, MaxMobs)
		decals  = make([]*Decal, MaxDecals)
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
	for i := 0; i < MaxDecals; i++ {
		decals[i] = NewDecal()
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
		Camera:           camera,
		Grid:             grid,
		State:            state,
		Mobs:             mobs,
		Decals:           decals,
		ActiveDecalCount: 0,
		ActiveMobCount:   0,
		entries:          entries,
		exit:             exit,
		blocks:           make(map[Ivec2]BlockPlacement),
		fearHistory:      fearHistory,
		fearIndex:        0,
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
			l.despawnMob(i)
		} else {
			mob.Update(elapsed, l)
		}
	}
}

func (l *Level) updateDecals(elapsed time.Duration) {
	for i := range l.Decals {
		decal := l.Decals[i]
		if decal.PendingDisable {
			l.disableDecal(i)
		} else {
			decal.Update(elapsed)
		}
	}
}

func (l *Level) updateSpawns(elapsed time.Duration) {
	// TODO: Calculate amount of charge as f(elapsed, rating)
	charge := 0.004 * math.Max(float64(l.State.Rating), 1)
	for i := range l.entries {
		entry := &l.entries[i]
		entry.AddCharge(charge)
		for entry.Spawn() {
			l.SpawnMob(entry.Pos)
		}
	}
}

func (l *Level) updateBlocks(elapsed time.Duration) {
	for pos, placement := range l.blocks {
		posV := mgl32.Vec2{float32(pos.X()), float32(pos.Y())}
		fear := placement.Block.FearPerSec * elapsed.Seconds()
		numHit := 0
		killed := make([]int, 0, placement.Block.MaxTargets)
		for i := range l.Mobs {
			mob := &l.Mobs[i]
			if numHit >= placement.Block.MaxTargets || !mob.Enabled {
				break
			}
			if mob.Pos.Sub(posV).Len() <= placement.Block.Range {
				numHit++
				if alive := mob.IncreaseFear(fear); !alive {
					// Mob has been scared to death.
					// TODO: uhhh this should be prettier.
					killed = append(killed, i)
					l.AddDecal(mob.Pos.Add(mgl32.Vec2{0, 0.5}), "ghost01_00", 2, 2*time.Second)
					l.State.Rating = l.penalizeRating()
				}
			}
		}
		if numHit > 0 {
			l.Grid.UpdateBlockState(pos, placement.Block, BlockScaring, placement.Variant)
		} else {
			l.Grid.UpdateBlockState(pos, placement.Block, BlockNormal, placement.Variant)
		}
		// Iterate from the back because we're doing some swapping and
		// don't wish to invalidate the rest of our indices.
		sort.Ints(killed)
		for i := len(killed) - 1; i > -1; i-- {
			l.disableMob(killed[i])
		}
	}
}

// Update computes a new simulation step for this level.
func (l *Level) Update(elapsed time.Duration) {
	l.updateBlocks(elapsed)
	l.updateMobs(elapsed)
	l.updateSpawns(elapsed)
	l.updateDecals(elapsed)
	l.Grid.Update(elapsed)
}

func (l *Level) SetMouse(screenX, screenY float32) {
	x, y := l.Camera.ScreenToWorldCoords(screenX, screenY)
	l.State.MousePos = mgl32.Vec2{x, y}
}

func (l *Level) GetMouse() mgl32.Vec2 {
	return l.State.MousePos
}

func (l *Level) SetCursor(frame string) {
	l.State.MouseCursor = frame
}

func (l *Level) GetCursor() string {
	return l.State.MouseCursor
}

func (l *Level) SetBlock(pos mgl32.Vec2, block *Block, variant int) {
	var (
		gridCoords = l.Grid.WorldToGrid(pos)
	)
	if center, ok := l.Grid.SetBlock(gridCoords, block, variant); ok {
		l.blocks[center] = BlockPlacement{block, variant}
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

func (l *Level) penalizeRating() int {
	sum := 0.0
	for i := range l.fearHistory {
		l.fearHistory[i] = math.Max(l.fearHistory[i]-1, 0)
		sum += l.fearHistory[i]
	}
	return int(math.Floor(sum/float64(len(l.fearHistory)) + 0.5))
}

func (l *Level) ClearHighlights() {
	l.Highlights = l.Highlights[0:0]
}

func (l *Level) SetHighlights(pos mgl32.Vec2, block *Block, variant int) {
	var (
		pre   = l.Grid.WorldToGrid(pos)
		post  = pre.Plus(block.Offset)
		frame = "special_squares_02"
	)
	l.ClearHighlights()
	if !l.Grid.IsBlockValid(pre, block, variant) {
		frame = "special_squares_03"
	}
	for y := 0; y < len(block.Variants[variant]); y++ {
		for x := 0; x < len(block.Variants[variant][y]); x++ {
			if block.Variants[variant][y][x] == nil {
				continue
			}
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

func (l *Level) despawnMob(i int) {
	var fear = l.Mobs[i].Fear
	switch {
	case fear < 5:
		l.AddDecal(l.Mobs[i].Pos.Add(mgl32.Vec2{0, 1.5}), "bubble_00", 1, 500*time.Millisecond)
	case fear > 8:
		l.AddDecal(l.Mobs[i].Pos.Add(mgl32.Vec2{0, 1.5}), "bubble_01", 1, 500*time.Millisecond)
	}
	l.fearHistory[l.fearIndex] = fear
	l.fearIndex = (l.fearIndex + 1) % len(l.fearHistory)
	l.State.Rating = l.calculateRating()
	l.State.Geld += int(math.Floor(fear*10.0 + 0.5))
	l.disableMob(i)
}

func (l *Level) disableMob(i int) {
	l.ActiveMobCount--
	l.Mobs[l.ActiveMobCount], l.Mobs[i] = l.Mobs[i], l.Mobs[l.ActiveMobCount]
	l.Mobs[l.ActiveMobCount].Disable()
}

func (l *Level) AddDecal(pos mgl32.Vec2, frame string, move float32, duration time.Duration) {
	if l.ActiveDecalCount >= MaxDecals {
		return
	}
	l.Decals[l.ActiveDecalCount].Activate(pos, frame, move, duration)
	l.ActiveDecalCount++
}

func (l *Level) disableDecal(i int) {
	if !l.Decals[i].Enabled {
		return
	}
	l.Decals[i].Disable()
	l.ActiveDecalCount--
	if l.ActiveDecalCount == i {
		return
	}
	l.Decals[l.ActiveDecalCount], l.Decals[i] = l.Decals[i], l.Decals[l.ActiveDecalCount]
}
