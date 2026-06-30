package main

import (
	"reflect"
	"sort"
	"testing"
)

func TestHeap(t *testing.T) {
	t.Run("basic push/pop", func(t *testing.T) {
		h := New()
		h.Push(3)
		h.Push(1)
		h.Push(2)

		got, ok := h.Pop()
		if !ok || got != 1 {
			t.Errorf("first pop = %d, ok=%v; want 1, true", got, ok)
		}

		got, ok = h.Pop()
		if !ok || got != 2 {
			t.Errorf("second pop = %d, ok=%v; want 2, true", got, ok)
		}

		got, ok = h.Pop()
		if !ok || got != 3 {
			t.Errorf("third pop = %d, ok=%v; want 3, true", got, ok)
		}

		_, ok = h.Pop()
		if ok {
			t.Error("pop from empty heap should return false")
		}
	})

	t.Run("popped order is non-decreasing", func(t *testing.T) {
		h := New()
		values := []int{5, 2, 8, 1, 9, 3}
		for _, v := range values {
			h.Push(v)
		}

		var popped []int
		for h.Len() > 0 {
			v, _ := h.Pop()
			popped = append(popped, v)
		}

		for i := 1; i < len(popped); i++ {
			if popped[i-1] > popped[i] {
				t.Errorf("popped sequence not non-decreasing: %v", popped)
			}
		}
	})

	t.Run("popped multiset equals pushed multiset", func(t *testing.T) {
		h := New()
		pushed := []int{5, 2, 8, 1, 9, 3, 5, 2}
		for _, v := range pushed {
			h.Push(v)
		}

		var popped []int
		for h.Len() > 0 {
			v, _ := h.Pop()
			popped = append(popped, v)
		}

		sort.Ints(pushed)
		sort.Ints(popped)
		if !reflect.DeepEqual(pushed, popped) {
			t.Errorf("popped multiset %v does not equal pushed multiset %v", popped, pushed)
		}
	})

	t.Run("empty heap edge cases", func(t *testing.T) {
		h := New()
		if h.Len() != 0 {
			t.Errorf("empty heap length = %d; want 0", h.Len())
		}

		_, ok := h.Pop()
		if ok {
			t.Error("pop from empty heap should return false")
		}
	})

	t.Run("single element", func(t *testing.T) {
		h := New()
		h.Push(42)
		if h.Len() != 1 {
			t.Errorf("heap with one element length = %d; want 1", h.Len())
		}

		got, ok := h.Pop()
		if !ok || got != 42 {
			t.Errorf("pop single element = %d, ok=%v; want 42, true", got, ok)
		}

		if h.Len() != 0 {
			t.Errorf("heap after popping single element length = %d; want 0", h.Len())
		}
	})
}

func FuzzHeap(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		h := New()
		var pushed []int
		for _, b := range data {
			x := int(int8(b)) // convert byte to signed int in range [-128, 127]
			h.Push(x)
			pushed = append(pushed, x)
		}

		var popped []int
		for h.Len() > 0 {
			v, _ := h.Pop()
			popped = append(popped, v)
		}

		// Check that popped sequence is non-decreasing
		for i := 1; i < len(popped); i++ {
			if popped[i-1] > popped[i] {
				t.Errorf("popped sequence not non-decreasing: %v", popped)
			}
		}

		// Check that popped multiset equals pushed multiset
		sort.Ints(pushed)
		sort.Ints(popped)
		if !reflect.DeepEqual(pushed, popped) {
			t.Errorf("popped multiset %v does not equal pushed multiset %v", popped, pushed)
		}
	})
}
