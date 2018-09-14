// Copyright © 2018 Chad Netzer <chad.netzer@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package inode

import (
	"reflect"
	"testing"
)

func TestSetAdd(t *testing.T) {
	var s Set
	s = NewSet()
	if len(s) != 0 {
		t.Errorf("Empty Set length isn't 0: %v", s)
	}
	s.Add(Ino(1))
	if len(s) != 1 {
		t.Errorf("Length %d Set should be 1", len(s))
	}
	s.Add(Ino(2))
	if len(s) != 2 {
		t.Errorf("Length %d Set should be 2", len(s))
	}
	s.Add(Ino(1))
	if len(s) != 2 {
		t.Errorf("Length %d Set after re-adding 1 should be 2", len(s))
	}
	s.Add(Ino(2))
	if len(s) != 2 {
		t.Errorf("Length %d Set after re-adding 2 should be 2", len(s))
	}
	s2 := NewSet(1, 2)
	if !reflect.DeepEqual(s, s2) {
		t.Errorf("Found unexpected unequal Sets: %v %v", s, s2)
	}
	s.Add(3)
	if reflect.DeepEqual(s, s2) {
		t.Errorf("Found unexpected equal Sets: %v %v", s, s2)
	}
}

func TestSetIntersection(t *testing.T) {
	s1 := NewSet(1)
	s2 := NewSet(2)

	// For each test, also test the symmetric operation
	in1 := s1.Intersection(s2)
	in2 := s2.Intersection(s1)
	if len(in1) != 0 || len(in2) != 0 {
		t.Errorf("Empty Set intersection length isn't 0: %v  %v", in1, in2)
	}

	s2.Add(Ino(1))
	in1 = s1.Intersection(s2)
	in2 = s2.Intersection(s1)
	if len(in1) != 1 || len(in2) != 1 {
		t.Errorf("Set intersection length isn't 1: %v  %v", in1, in2)
	}
	if !reflect.DeepEqual(in1, NewSet(Ino(1))) {
		t.Errorf("Set intersection doesn't contain only 1: %v", in1)
	}
	if !reflect.DeepEqual(in2, NewSet(Ino(1))) {
		t.Errorf("Set intersection doesn't contain only 1: %v", in2)
	}

	s1.Add(Ino(2))
	in1 = s1.Intersection(s2)
	in2 = s2.Intersection(s1)
	if len(in1) != 2 || len(in2) != 2 {
		t.Errorf("Set intersection length isn't 2: %v  %v", in1, in2)
	}
	if !reflect.DeepEqual(in1, NewSet(Ino(1), Ino(2))) {
		t.Errorf("Set intersection isn't {1,2}: %v", in1)
	}
	if !reflect.DeepEqual(in2, NewSet(Ino(1), Ino(2))) {
		t.Errorf("Set intersection isn't {1,2}: %v", in2)
	}

	s1 = NewSet(Ino(1))
	s2 = NewSet(Ino(2))
	s3 := NewSet(Ino(3))
	in1 = SetIntersection(s1, s2, s3)
	if !reflect.DeepEqual(in1, NewSet()) {
		t.Errorf("Set intersection isn't empty: %v", in1)
	}

	s1 = NewSet(Ino(1))
	s2 = NewSet(Ino(1))
	s3 = NewSet(Ino(1))
	in1 = SetIntersection(s1, s2, s3)
	if !reflect.DeepEqual(in1, NewSet(Ino(1))) {
		t.Errorf("Set intersection isn't {1}: %v", in1)
	}

	s1 = NewSet(Ino(1), Ino(2), Ino(4))
	s2 = NewSet(Ino(2), Ino(3), Ino(4))
	s3 = NewSet(Ino(1), Ino(3), Ino(4))
	in1 = SetIntersection(s1, s2, s3)
	if !reflect.DeepEqual(in1, NewSet(Ino(4))) {
		t.Errorf("Set intersection isn't {4}: %v", in1)
	}
}

func TestSetDifference(t *testing.T) {
	s1 := NewSet()
	s2 := NewSet()
	diff1 := s1.Difference(s2)
	diff2 := s2.Difference(s1)
	if len(diff1) != 0 || len(diff2) != 0 {
		t.Errorf("Empty Set difference length isn't 0: %v  %v", diff1, diff2)
	}

	s1.Add(Ino(1))
	diff1 = s1.Difference(s2)
	if !reflect.DeepEqual(diff1, NewSet(Ino(1))) {
		t.Errorf("Set difference doesn't contain only 1: %v", diff1)
	}
	diff2 = s2.Difference(s1)
	if !reflect.DeepEqual(diff2, NewSet()) {
		t.Errorf("Set difference isn't empty: %v", diff2)
	}

	s2.Add(Ino(1))
	diff1 = s1.Difference(s2)
	if !reflect.DeepEqual(diff1, NewSet()) {
		t.Errorf("Set difference isn't empty: %v", diff1)
	}
	diff2 = s2.Difference(s1)
	if !reflect.DeepEqual(diff2, NewSet()) {
		t.Errorf("Set difference isn't empty: %v", diff2)
	}

	s1.Add(Ino(2))
	diff1 = s1.Difference(s2)
	if !reflect.DeepEqual(diff1, NewSet(Ino(2))) {
		t.Errorf("Set difference doesn't contain only 2: %v", diff1)
	}
	diff2 = s2.Difference(s1)
	if !reflect.DeepEqual(diff2, NewSet()) {
		t.Errorf("Set difference isn't empty: %v", diff2)
	}
}

func TestSetAsSlice(t *testing.T) {
	s1 := NewSet()
	if !reflect.DeepEqual(s1.AsSlice(), []Ino{}) {
		t.Errorf("Set.AsSlice() isn't empty: %v", s1.AsSlice())
	}
	s1.Add(Ino(1))
	if !reflect.DeepEqual(s1.AsSlice(), []Ino{Ino(1)}) {
		t.Errorf("Set.AsSlice() isn't [1]: %v", s1.AsSlice())
	}
	s1.Add(Ino(2))
	if len(s1) != 2 {
		t.Errorf("Length of Set.AsSlice() isn't 2: %v", s1.AsSlice())
	}
	if len(append(NewSet().AsSlice(), NewSet().AsSlice()...)) != 0 {
		t.Errorf("Length of appended Set.AsSlice() was expected to be 0")
	}
	if len(append(NewSet(Ino(1)).AsSlice(), NewSet().AsSlice()...)) != 1 {
		t.Errorf("Length of appended Set.AsSlice() was expected to be 1")
	}
	if len(append(NewSet(Ino(1)).AsSlice(), NewSet(Ino(2)).AsSlice()...)) != 2 {
		t.Errorf("Length of appended Set.AsSlice() was expected to be 2")
	}
}
