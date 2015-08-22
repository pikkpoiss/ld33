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

type Level struct {
	Grid *Grid
	Mobs []*Mob
}

func NewLevel() (level *Level, err error) {
	level = &Level{
		Grid: NewGrid(),
		Mobs: []*Mob{},
	}
	return
}

func (l *Level) Update(elapsed time.Duration) {
	for _, mob := range l.Mobs {
		mob.Update(elapsed, l)
	}
}

func (l *Level) AddMob(pos mgl32.Vec2) {
	l.Mobs = append(l.Mobs, &Mob{pos, 2.0})
}
