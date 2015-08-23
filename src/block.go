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

type BlockPlacement struct {
	Pos Ivec2
	Block   *Block
	Variant int
}

type BlockAnimations map[BlockState][]int

var (
	SkeletonAnimations = BlockAnimations{
		BlockNormal:  []int{0},
		BlockScaring: []int{1, 2, 3, 4},
	}
	SkeletonTemplate = &GridItemTemplate{
		false,
		"skeleton01_%02v",
		SkeletonAnimations,
	}
)

var (
	SpikesAnimations = BlockAnimations{
		BlockNormal: []int{0},
		BlockScaring: []int{
			1,
			2, 2, 2, 2, 2, 2, 2, 2, 2,
			3, 3,
			4, 4,
			0, 0, 0, 0, 0, 0, 0, 0,
		},
	}
	SpikesTemplate = &GridItemTemplate{
		false,
		"spikes01_%02v",
		SpikesAnimations,
	}
)

type GridItemTemplate struct {
	Passable bool
	Frame    string
	Frames   BlockAnimations
}

type BlockTemplate [][]*GridItemTemplate

// TODO: Introduce a cooldown for scaring people.
type Block struct {
	Variants   []BlockTemplate
	Offset     Ivec2
	Range      float32 // Radius of effectiveness.
	MaxTargets int     // -1 for infinite.
	FearPerSec float64 // Amount of fear added to target per second.
	Cost       int
	Title      string
}

var (
	SkellyBlock = Block{
		Variants: []BlockTemplate{
			BlockTemplate{
				[]*GridItemTemplate{
					SkeletonTemplate,
				},
			},
		},
		Offset:     Ivec2{0, 0},
		Range:      1.5,
		MaxTargets: 1,
		FearPerSec: 2.0,
		Cost:       10,
		Title:      "Mr. Bones",
	}

	SpikesBlock = Block{
		Variants: []BlockTemplate{
			BlockTemplate{
				[]*GridItemTemplate{SpikesTemplate, SpikesTemplate, SpikesTemplate},
				[]*GridItemTemplate{nil, nil, nil},
				[]*GridItemTemplate{SpikesTemplate, SpikesTemplate, SpikesTemplate},
			},
			BlockTemplate{
				[]*GridItemTemplate{SpikesTemplate, nil, SpikesTemplate},
				[]*GridItemTemplate{SpikesTemplate, nil, SpikesTemplate},
				[]*GridItemTemplate{SpikesTemplate, nil, SpikesTemplate},
			},
		},
		Offset:     Ivec2{-1, -1},
		Range:      5.0,
		MaxTargets: 3,
		FearPerSec: 0.5,
		Cost:       100,
		Title:      "Spiketron 5000",
	}

	CornerBlock = Block{
		Variants: []BlockTemplate{
			BlockTemplate{
				[]*GridItemTemplate{SpikesTemplate, SpikesTemplate, SpikesTemplate},
				[]*GridItemTemplate{nil, nil, SpikesTemplate},
				[]*GridItemTemplate{SpikesTemplate, nil, SpikesTemplate},
			},
			BlockTemplate{
				[]*GridItemTemplate{SpikesTemplate, nil, SpikesTemplate},
				[]*GridItemTemplate{nil, nil, SpikesTemplate},
				[]*GridItemTemplate{SpikesTemplate, SpikesTemplate, SpikesTemplate},
			},
			BlockTemplate{
				[]*GridItemTemplate{SpikesTemplate, nil, SpikesTemplate},
				[]*GridItemTemplate{SpikesTemplate, nil, nil},
				[]*GridItemTemplate{SpikesTemplate, SpikesTemplate, SpikesTemplate},
			},
			BlockTemplate{
				[]*GridItemTemplate{SpikesTemplate, SpikesTemplate, SpikesTemplate},
				[]*GridItemTemplate{SpikesTemplate, nil, nil},
				[]*GridItemTemplate{SpikesTemplate, nil, SpikesTemplate},
			},
		},
		Offset:     Ivec2{-1, -1},
		Range:      5.0,
		MaxTargets: 3,
		FearPerSec: 0.5,
		Cost:       100,
		Title:      "Spiketron 6000 GT",
	}
)
