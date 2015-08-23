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
)

type Ivec2 [2]int32

func (i Ivec2) X() int32 {
	return i[0]
}

func (i Ivec2) Y() int32 {
	return i[1]
}

func (i Ivec2) Plus(a Ivec2) Ivec2 {
	return Ivec2{i[0] + a[0], i[1] + a[1]}
}

type Grid struct {
	background *twodee.Grid
	grid       *twodee.Grid
	sources    []Ivec2
	sink       Ivec2
}

func NewGrid() (g *Grid, err error) {
	var (
		background *twodee.Grid
		grid       *twodee.Grid
	)
	if background, err = LoadTiledMap("resources/maps/map01.tmx"); err != nil {
		return
	}
	grid = twodee.NewGrid(background.Width, background.Height, 1.0)
	g = &Grid{
		background: background,
		grid:       grid,
	}
	return
}

func (g *Grid) AddSource(pt Ivec2) {
	g.Set(pt, &GridItem{false, 0, "special_squares_00"})
	g.sources = append(g.sources, pt)
}

func (g *Grid) SetSink(pt Ivec2) {
	g.Set(pt, &GridItem{false, 0, "special_squares_00"})
	g.sink = pt
}

func (g *Grid) Set(pt Ivec2, item *GridItem) {
	g.grid.Set(pt.X(), pt.Y(), item)
}

func (g *Grid) IsBlockValid(origin Ivec2, block *Block) (ok bool) {
	var (
		pt   = origin.Plus(block.Offset)
		item *GridItem
	)
	ok = true
	for y := 0; y < len(block.Template); y++ {
		for x := 0; x < len(block.Template[y]); x++ {
			item = g.Get(pt.Plus(Ivec2{int32(x), int32(y)}))
			if item != nil && !item.Passable() {
				ok = false
				break
			}
		}
	}
	return
}

// SetBlock attempts to place the block in a "centered" fashion on the given
// origin. It returns the calculated center as well as a bool indicating
// whether placement was successful.
func (g *Grid) SetBlock(origin Ivec2, block *Block) (Ivec2, bool) {
	var (
		pt = origin.Plus(block.Offset)
	)
	if !g.IsBlockValid(origin, block) {
		return pt, false
	}
	for y := 0; y < len(block.Template); y++ {
		for x := 0; x < len(block.Template[y]); x++ {
			g.Set(pt.Plus(Ivec2{int32(x), int32(y)}), block.Template[y][x])
		}
	}
	return pt, true
}

func (g *Grid) Width() int32 {
	return g.grid.Width
}

func (g *Grid) Height() int32 {
	return g.grid.Height
}

func (g *Grid) WorldToGrid(worldCoords mgl32.Vec2) Ivec2 {
	return Ivec2{
		g.grid.GridPosition(worldCoords[0]),
		g.grid.GridPosition(worldCoords[1]),
	}
}

func (g *Grid) GetNextStepToSink(pt mgl32.Vec2) (out mgl32.Vec2, dist int32, valid bool) {
	var (
		gridPt = g.WorldToGrid(pt)
		points []Ivec2
		item   *GridItem
	)
	dist = 9999999
	points = g.getAdjacent(gridPt)
	valid = false
	for _, adj := range points {
		item = g.Get(adj)
		if item == nil {
			item = g.GetBg(adj)
		}
		if item != nil && item.Passable() {
			if item.Distance() < dist {
				dist = item.Distance()
				out = mgl32.Vec2{
					g.grid.InversePosition(adj.X()),
					g.grid.InversePosition(adj.Y()),
				}
				valid = true
			}
		}
	}
	return
}

func (g *Grid) resetDistances() {
	var (
		x    int32
		y    int32
		item *GridItem
		pt   Ivec2
	)
	for x = 0; x < g.Width(); x++ {
		for y = 0; y < g.Height(); y++ {
			pt = Ivec2{x, y}
			item = g.Get(pt)
			if item == nil {
				item = g.GetBg(pt)
			}
			if item != nil && item.Passable() {
				item.SetDistance(-1)
			}
		}
	}
}

func (g *Grid) CalculateDistances() {
	var (
		queue       = []Ivec2{g.sink}
		dist  int32 = 1
		item  *GridItem
	)
	g.resetDistances()
	for len(queue) > 0 {
		pt := queue[0]
		queue = queue[1:]
		points := g.getAdjacent(pt)
		item = g.Get(pt)
		if item == nil {
			item = g.GetBg(pt)
		}
		dist = item.Distance() + 1
		for _, adj := range points {
			item = g.Get(adj)
			if item == nil {
				item = g.GetBg(adj)
			}
			if item != nil && item.Distance() == -1 && item.Passable() {
				queue = append(queue, adj)
				item.SetDistance(dist)
			}
		}
	}
}

func (g *Grid) Get(pt Ivec2) (item *GridItem) {
	var (
		tditem twodee.GridItem
		ok     bool
	)
	tditem = g.grid.Get(pt[0], pt[1])
	if item, ok = tditem.(*GridItem); ok {
		return item
	}
	return nil
}

func (g *Grid) GetBg(pt Ivec2) (item *GridItem) {
	var (
		tditem twodee.GridItem
		ok     bool
	)
	tditem = g.background.Get(pt[0], pt[1])
	if item, ok = tditem.(*GridItem); ok {
		return item
	}
	return nil
}

func (g *Grid) getAdjacent(point Ivec2) (points []Ivec2) {
	var adj Ivec2
	adj = point.Plus(Ivec2{-1, 0})
	if adj.X() < g.Width() && adj.X() >= 0 {
		points = append(points, adj)
	}
	adj = point.Plus(Ivec2{1, 0})
	if adj.X() < g.Width() && adj.X() >= 0 {
		points = append(points, adj)
	}
	adj = point.Plus(Ivec2{0, -1})
	if adj.Y() < g.Height() && adj.Y() >= 0 {
		points = append(points, adj)
	}
	adj = point.Plus(Ivec2{0, 1})
	if adj.Y() < g.Height() && adj.Y() >= 0 {
		points = append(points, adj)
	}
	return points
}
