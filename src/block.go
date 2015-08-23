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

type BlockState int32

const (
	_                      = iota
	BlockNormal BlockState = 1 << iota
	BlockScaring
)

type BlockAnimations map[BlockState][]int

var (
	SkeletonAnimations = BlockAnimations{
		BlockNormal:  []int{0},
		BlockScaring: []int{1, 2, 3, 4},
	}
	SkeletonTemplate = &BlockTemplate{
		false,
		"skeleton01_%02v",
		SkeletonAnimations,
	}
)

var (
	SpikesAnimations = BlockAnimations{
		BlockNormal:  []int{0},
		BlockScaring: []int{1, 2, 3, 3, 3, 3, 3, 0, 0, 0, 0},
	}
	SpikesTemplate = &BlockTemplate{
		false,
		"spikes01_%02v",
		SpikesAnimations,
	}
)

type BlockTemplate struct {
	Passable bool
	Frame    string
	Frames   BlockAnimations
}

// TODO: Introduce a cooldown for scaring people.
type Block struct {
	Template   [][]*BlockTemplate
	Offset     Ivec2
	Range      float32 // Radius of effectiveness.
	MaxTargets int     // -1 for infinite.
	FearPerSec float64 // Amount of fear added to target per second.
	Cost       int
}

var (
	OneBlock = Block{
		Template: [][]*BlockTemplate{
			[]*BlockTemplate{
				SkeletonTemplate,
			},
		},
		Offset:     Ivec2{0, 0},
		Range:      1.5,
		MaxTargets: 1,
		FearPerSec: 2.0,
		Cost:       10,
	}

	ThreeBlock = Block{
		Template: [][]*BlockTemplate{
			[]*BlockTemplate{
				SpikesTemplate,
				SpikesTemplate,
				SpikesTemplate,
			},
			[]*BlockTemplate{
				nil,
				nil,
				nil,
			},
			[]*BlockTemplate{
				SpikesTemplate,
				SpikesTemplate,
				SpikesTemplate,
			},
		},
		Offset:     Ivec2{-1, -1},
		Range:      5.0,
		MaxTargets: 3,
		FearPerSec: 0.5,
		Cost:       100,
	}
)
