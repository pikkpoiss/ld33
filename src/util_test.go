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
	"testing"
)

var bufferTests = []struct {
	entries  []float64
	expected float64
}{
	{[]float64{0}, 0.0},
	{[]float64{0, 1}, 0.5},
	{[]float64{0, 1, 2}, 1.0},
	{[]float64{0, 1, 2, 3}, 6.0 / 4.0},
	{[]float64{0, 1, 2, 3, 4, 5}, 15.0 / 5.0},
}

func TestCircularBuffer(t *testing.T) {
	b := NewCircularBuffer(5)
	eSample := 0.0
	if sample := b.Sample(); sample != eSample {
		t.Fatalf("Expected %v average got %v", eSample, sample)
	}

	for _, bt := range bufferTests {
		b = NewCircularBuffer(5)
		for _, e := range bt.entries {
			b.AddEntry(e)
		}
		if sample := b.Sample(); sample != bt.expected {
			t.Fatalf("Expected %v average got %v", bt.expected, sample)
		}
	}
}
