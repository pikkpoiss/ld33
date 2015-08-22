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
	grid        *twodee.Grid
	defaultItem *GridItem
	enter       Ivec2
	exit        Ivec2
}

func NewGrid() (g *Grid) {
	g = &Grid{
		grid:        twodee.NewGrid(64, 40, PxPerUnit),
		defaultItem: &GridItem{true, 0},
	}
	g.SetEnter(4, 19)
	g.SetExit(60, 20)
	g.init()
	g.setDistances()
	return
}

func (g *Grid) SetEnter(x, y int32) {
	g.grid.Set(x, y, &GridItem{false, 0})
	g.enter = Ivec2{x, y}
}

func (g *Grid) SetExit(x, y int32) {
	g.grid.Set(x, y, &GridItem{false, 0})
	g.exit = Ivec2{x, y}
}

func (g *Grid) Get(x, y int32) (item *GridItem) {
	item = g.get(Ivec2{x, y})
	if item == nil {
		return g.defaultItem
	}
	return item
}

func (g *Grid) Width() int32 {
	return g.grid.Width
}

func (g *Grid) Height() int32 {
	return g.grid.Height
}

func (g *Grid) init() {
	var (
		x    int32
		y    int32
		item *GridItem
	)
	for x = 0; x < g.Width(); x++ {
		for y = 0; y < g.Height(); y++ {
			item = g.get(Ivec2{x, y})
			if item == nil {
				g.grid.Set(x, y, &GridItem{false, -1})
			}
		}
	}
}

func (g *Grid) setDistances() {
	var (
		queue       = []Ivec2{g.exit}
		dist  int32 = 1
		item  *GridItem
	)
	for len(queue) > 0 {
		pt := queue[0]
		queue = queue[1:]
		points := g.getAdjacent(pt)
		item = g.get(pt)
		dist = item.Distance() + 1
		for _, adj := range points {
			item = g.get(adj)
			if item != nil && item.Distance() == -1 {
				queue = append(queue, adj)
				item.SetDistance(dist)
			}
		}
	}
}

func (g *Grid) get(pt Ivec2) (item *GridItem) {
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
