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

type Grid struct {
	grid        *twodee.Grid
	defaultItem *GridItem
}

func NewGrid() (g *Grid) {
	g = &Grid{
		grid:        twodee.NewGrid(64, 40, PxPerUnit),
		defaultItem: &GridItem{true},
	}
	g.grid.Set(4, 19, &GridItem{false})
	g.grid.Set(60, 20, &GridItem{false})
	return
}

func (g *Grid) Get(x, y int32) (item *GridItem) {
	var (
		tditem twodee.GridItem
		ok     bool
	)
	tditem = g.grid.Get(x, y)
	if item, ok = tditem.(*GridItem); !ok {
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
