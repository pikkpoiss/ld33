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

type CircularBuffer struct {
	buffer                []float64
	idx, numEntries, size int
	sum                   float64
}

func NewCircularBuffer(size int) *CircularBuffer {
	return &CircularBuffer{
		buffer:     make([]float64, size),
		idx:        0,
		numEntries: 0,
		size:       size,
		sum:        0,
	}
}

func (b *CircularBuffer) AddEntry(v float64) {
	pVal := b.buffer[b.idx]
	b.sum -= pVal
	b.buffer[b.idx] = v
	b.sum += v
	b.idx = (b.idx + 1) % b.size
	b.numEntries = minInt(b.size, b.numEntries+1)
}

// Sample returns the average of all entries collected by the circular buffer.
func (b *CircularBuffer) Sample() float64 {
	if b.numEntries == 0 {
		return 0.0
	}
	return b.sum / float64(b.numEntries)
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
