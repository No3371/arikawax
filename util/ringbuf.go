package util

import "iter"

// RingBuffer is a circular buffer that can be used to store unlimited number of items.
// It automatically expands when the buffer is full.
type RingBuffer[T comparable] struct {
	buffer []T
	read   int
	write  int
	size   int
}

func (r *RingBuffer[T]) Init(size int) {
	size += 1
	r.buffer = make([]T, size)
	r.size = size
	r.read = 0
	r.write = 0
}

func (r *RingBuffer[T]) IsFull() bool {
	return (r.write+1)%r.size == r.read
}

func (r *RingBuffer[T]) IsEmpty() bool {
	return r.read == r.write
}

func (r *RingBuffer[T]) Push(item T) {
	r.buffer[r.write] = item
	if r.IsFull() {
		r.size *= 2
		buf := make([]T, r.size)
		copy(buf, r.buffer)
		r.buffer = buf
	}
	r.write = (r.write + 1) % r.size
}

func (r *RingBuffer[T]) PushUnique(item T) {
	collide := false
	for v := range r.PeekAll() {
		if v == item {
			collide = true
			break
		}
	}
	if !collide {
		r.Push(item)
	}
}

func (r *RingBuffer[T]) Len() int {
	if r.IsEmpty() {
		return 0
	}
	if r.read < r.write {
		return r.write - r.read
	}
	return r.size - r.read + r.write
}

func (r *RingBuffer[T]) Pop() (v T, ok bool) {
	if r.buffer == nil {
		return
	}
	if r.IsEmpty() { // exhausted
		return
	}
	v = r.buffer[r.read]
	var empty T
	r.buffer[r.read] = empty
	r.read = (r.read + 1) % r.size
	return v, true
}

func (r *RingBuffer[T]) Peek() (v T, ok bool) {
	if r.IsEmpty() {
		return
	}
	v = r.buffer[r.read]
	return v, true
}

func (r *RingBuffer[T]) PeekAll() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := r.read; i != r.write; i = (i + 1) % r.size {
			if !yield(r.buffer[i]) {
				return
			}
		}
	}
}

func (r *RingBuffer[T]) PeekAllReverse() iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := (r.write - 1 + r.size) % r.size; i != (r.read-1+r.size)%r.size; i = (i - 1 + r.size) % r.size {
			if !yield(r.buffer[i]) {
				return
			}
		}
	}
}
