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

const (
	FAIL_RATING  = 0
	WIN_RATING   = 9
	WIN_DURATION = 5 * time.Second
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
	App              *Application
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
	highlighted      *BlockPlacement
	gameEventHandler *twodee.GameEventHandler
	durAtWinRating   time.Duration
}

const (
	MaxMobs   = 200
	MaxDecals = 10
)

func NewLevel(state *State, sheet *twodee.Spritesheet, gameEventHandler *twodee.GameEventHandler) (level *Level, err error) {
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
		exit        = NewSpawnZone(Ivec2{24, 9})
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
		gameEventHandler: gameEventHandler,
		durAtWinRating:   0,
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
					l.gameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayDeathEffect))
					l.State.Rating = l.penalizeRating()
				}
			}
		}
		if numHit > 0 {
			l.Grid.UpdateBlockState(pos, placement.Block, BlockScaring, placement.Variant)
			switch placement.Block.Title {
			case "Mr. Bones":
				l.gameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayMrBonesEffect))
			case "Spiketron 5000":
				l.gameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlaySpikesEffect))
			case "Spiketron 6000 GT":
				l.gameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlaySpikesEffect))
			}
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

// checkConditions checks to see if the player has lost. If so, it enqueues a
// PlayerLost event.
func (l *Level) checkConditions(elapsed time.Duration) {
	if l.State.Rating <= FAIL_RATING {
		l.gameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayerLost))
	}
	if l.State.Rating >= WIN_RATING {
		l.durAtWinRating += elapsed
		if l.durAtWinRating >= WIN_DURATION {
			l.gameEventHandler.Enqueue(twodee.NewBasicGameEvent(PlayerWon))
		}
	} else {
		l.durAtWinRating = 0
	}
}

// Update computes a new simulation step for this level.
func (l *Level) Update(elapsed time.Duration) {
	l.updateBlocks(elapsed)
	l.updateMobs(elapsed)
	l.updateSpawns(elapsed)
	l.updateDecals(elapsed)
	l.Grid.Update(elapsed)
	l.checkConditions(elapsed)
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
		placement  = BlockPlacement{gridCoords, block, variant}
	)
	if center, ok := l.Grid.SetBlock(placement); ok {
		l.blocks[center] = placement
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

func (l *Level) clearHighlights() {
	l.Highlights = l.Highlights[0:0]
}

func (l *Level) SetHighlights(pos mgl32.Vec2, block *Block, variant int) {
	l.highlighted = &BlockPlacement{
		Pos:     l.Grid.WorldToGrid(pos),
		Block:   block,
		Variant: variant,
	}
	l.RefreshHighlights()
}

func (l *Level) UnsetHighlights() {
	l.clearHighlights()
	l.highlighted = nil
}

func (l *Level) RefreshHighlights() {
	if l.highlighted == nil {
		return
	}
	var (
		pre   = l.highlighted.Pos
		post  = pre.Plus(l.highlighted.Block.Offset)
		frame = "special_squares_02"
	)
	l.clearHighlights()
	if !l.Grid.IsBlockValid(*l.highlighted) || l.State.Geld < l.highlighted.Block.Cost {
		frame = "special_squares_03"
	}
	for y := 0; y < len(l.highlighted.Block.Variants[l.highlighted.Variant]); y++ {
		for x := 0; x < len(l.highlighted.Block.Variants[l.highlighted.Variant][y]); x++ {
			if l.highlighted.Block.Variants[l.highlighted.Variant][y][x] == nil {
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

func (l *Level) AddGeld(amount int) {
	l.State.Geld += amount
	l.RefreshHighlights()
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
	l.AddGeld(int(math.Floor(fear*10.0 + 0.5)))
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
