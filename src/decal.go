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

type Decal struct {
	Pos            mgl32.Vec2
	Adjust         float32
	Enabled        bool
	PendingDisable bool
	Frame          string
	Animation      twodee.Animator
}

func NewDecal() *Decal {
	return &Decal{
		Enabled:        false,
		PendingDisable: false,
	}
}

func (d *Decal) Activate(pos mgl32.Vec2, frame string, move float32, duration time.Duration) {
	d.Pos = pos
	d.Frame = frame
	d.Adjust = 0
	d.Enabled = true
	d.PendingDisable = false
	d.Animation = twodee.NewEaseOutAnimation(&d.Adjust, 0, move, duration)
	d.Animation.SetCallback(func() {
		d.PendingDisable = true
	})
	return
}

func (d *Decal) Disable() {
	d.Enabled = false
	d.PendingDisable = false
}

func (d *Decal) Update(elapsed time.Duration) {
	if d.Animation != nil {
		d.Animation.Update(elapsed)
	}
}

func (d *Decal) SpriteConfig(sheet *twodee.Spritesheet) twodee.SpriteConfig {
	var (
		frame = sheet.GetFrame(d.Frame)
	)
	return twodee.SpriteConfig{
		View: twodee.ModelViewConfig{
			d.Pos.X(), d.Pos.Y() + d.Adjust, 0.0,
			0, 0, 0,
			1.0, 1.0, 1.0,
		},
		Frame: frame.Frame,
	}
}
