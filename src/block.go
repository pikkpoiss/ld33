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

import ()

// TODO: Introduce a cooldown for scaring people.
type Block struct {
	Template   [][]*GridItem
	Offset     Ivec2
	Range      float32 // Radius of effectiveness.
	MaxTargets int     // -1 for infinite.
	FearPerSec float64 // Amount of fear added to target per second.
	Cost       int
}

var (
	OneBlock = Block{
		[][]*GridItem{
			[]*GridItem{
				NewGridItem(false, "skeleton01_00"),
			},
		},
		Ivec2{0, 0},
		1.5,
		1,
		4.0,
		10,
	}

	ThreeBlock = Block{
		[][]*GridItem{
			[]*GridItem{
				NewGridItem(false, "spikes01_00"),
				NewGridItem(false, "spikes01_00"),
				NewGridItem(false, "spikes01_00"),
			},
			[]*GridItem{
				nil,
				nil,
				nil,
			},
			[]*GridItem{
				NewGridItem(false, "spikes01_00"),
				NewGridItem(false, "spikes01_00"),
				NewGridItem(false, "spikes01_00"),
			},
		},
		Ivec2{-1, -1},
		5.0,
		3,
		0,
		100,
	}
)
