package util

import (
	"testing"
)

func TestRingBufferPush(t *testing.T) {
	rb := RingBuffer[int]{}
	rb.Init(5)
	rb.Push(1)
	rb.Push(2)
	rb.Push(3)
	rb.Push(4)
	rb.Push(5)
	rb.Push(6)
	if rb.Len() != 6 {
		t.Errorf("Expected length 6, got %d", rb.Len())
	}
	if v, _ := rb.Pop(); v != 1 {
		t.Errorf("Expected 1, got %d", v)
	}
}

func TestRingBufferLen(t *testing.T) {
	rb := RingBuffer[int]{}
	rb.Init(5)
	if rb.Len() != 0 {
		t.Errorf("Expected length 0, got %d", rb.Len())
	}

	rb.Push(1)
	rb.Push(2)
	rb.Push(3)
	if rb.Len() != 3 {
		t.Errorf("Expected length 3, got %d", rb.Len())
	}
	rb.Push(4)
	rb.Push(5)
	rb.Push(6)
	if rb.Len() != 6 {
		t.Errorf("Expected length 6, got %d", rb.Len())
	}

	rb.Pop()
	rb.Pop()
	rb.Pop()
	if rb.Len() != 3 {
		t.Errorf("Expected length 3, got %d", rb.Len())
	}

	rb.Pop()
	rb.Pop()
	rb.Pop()

	if rb.Len() != 0 {
		t.Errorf("Expected length 0, got %d", rb.Len())
	}
}

func TestRingBufferPushUnique(t *testing.T) {
	rb := RingBuffer[int]{}
	rb.Init(5)
	rb.PushUnique(1)
	rb.PushUnique(2)
	rb.PushUnique(3)
	rb.PushUnique(3)
	rb.PushUnique(3)
	rb.PushUnique(4)

	count3 := 0
	for _, v := range rb.buffer {
		if v == 3 {
			count3++
		}
	}
	if count3 != 1 {
		t.Errorf("Expected 1, got %d", count3)
	}
}

func TestRingBufferPop(t *testing.T) {
	rb := RingBuffer[int]{}
	rb.Init(5)

	rb.Push(1)
	rb.Push(2)
	rb.Push(3)
	rb.Push(4)

	if rb.Len() != 4 {
		t.Errorf("Expected length 4, got %d", rb.Len())
	}

	rb.Push(5)

	if rb.Len() != 5 {
		t.Errorf("Expected length 5, got %d", rb.Len())
	}

	popped := 0
	for i := 0; i < 5; i++ {
		v, ok := rb.Pop()
		if ok {
			t.Logf("Popped: %d, %t", v, ok)
			popped++
		}
	}
	if rb.Len() != 0 {
		t.Errorf("Expected length 0, got %d", rb.Len())
	}
	if popped != 5 {
		t.Errorf("Expected 5, got %d", popped)
	}
}

func TestRingBufferPop2(t *testing.T) {
	rb := RingBuffer[int]{}
	rb.Init(5)

	rb.Push(1)
	rb.Push(2)
	rb.Push(3)
	rb.Push(4)
	rb.Push(5)
	rb.Push(6)

	if rb.Len() != 6 {
		t.Errorf("Expected length 6, got %d", rb.Len())
	}

	popped := 0
	for i := 0; i < 7; i++ {
		v, ok := rb.Pop()
		if ok {
			t.Logf("Popped: %d, %t", v, ok)
			popped++
		}
	}
	if rb.Len() != 0 {
		t.Errorf("Expected length 0, got %d", rb.Len())
	}
	if popped != 6 {
		t.Errorf("Expected popped 6, got %d", popped)
	}
}

func TestRingBufferPeekAll(t *testing.T) {
	rb := RingBuffer[int]{}
	rb.Init(5)

	rb.Push(1)
	rb.Push(2)
	rb.Push(3)
	rb.Push(4)
	rb.Push(5)
	rb.Push(6)

	rb.Pop()
	rb.Pop()

	peeked := []int{}
	for v := range rb.PeekAll() {
		peeked = append(peeked, v)
	}
	if len(peeked) != 4 {
		t.Errorf("Expected 4, got %d", len(peeked))
	}
	expected := []int{3, 4, 5, 6}
	for i, v := range peeked {
		if v != expected[i] {
			t.Errorf("Expected %d, got %d", expected[i], v)
		}
	}
}

func TestRingBufferPeekAllReverse(t *testing.T) {
	rb := RingBuffer[int]{}
	rb.Init(5)

	rb.Push(1)
	rb.Push(2)
	rb.Push(3)
	rb.Push(4)
	rb.Push(5)
	rb.Push(6)

	rb.Pop()
	rb.Pop()

	peeked := []int{}
	for v := range rb.PeekAllReverse() {
		peeked = append(peeked, v)
	}
	if len(peeked) != 4 {
		t.Errorf("Expected peeked 4, got %d", len(peeked))
	}
	expected := []int{6, 5, 4, 3}
	for i, v := range peeked {
		if v != expected[i] {
			t.Errorf("Expected %d, got %d", expected[i], v)
		}
	}
}
