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

type GridItem struct {
	passable bool
	distance int32
	Frame    string
}

func (i *GridItem) Passable() bool {
	return i.passable
}

func (i *GridItem) Opaque() bool {
	return i.passable
}

func (i *GridItem) SetDistance(dist int32) {
	i.distance = dist
}

func (i *GridItem) Distance() int32 {
	return i.distance
}
