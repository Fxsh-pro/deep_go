package main

import "errors"

type CircularQueue struct {
	values []int
	length int
	head   int
	tail   int
}

func NewCircularQueue(size int) (CircularQueue, error) {
	if size <= 0 {
		return CircularQueue{}, errors.New("size must be greater than 0")
	}
	return CircularQueue{
		values: make([]int, size),
		length: size,
		head:   -1,
		tail:   -1,
	}, nil
}

func (q *CircularQueue) Push(value int) bool {
	if q.Full() {
		return false
	}
	if q.head == -1 {
		q.head = 0
	}
	q.tail = (q.tail + 1) % q.length
	q.values[q.tail] = value
	return true
}

func (q *CircularQueue) Pop() bool {
	if q.Empty() {
		return false
	}
	if q.head == q.tail {
		q.head, q.tail = -1, -1
		return true
	}
	q.head = (q.head + 1) % q.length
	return true
}

func (q *CircularQueue) Front() int {
	if q.Empty() {
		return -1
	}
	return q.values[q.head]
}

func (q *CircularQueue) Back() int {
	if q.Empty() {
		return -1
	}
	return q.values[q.tail]
}

func (q *CircularQueue) Empty() bool {
	return q.head == -1
}

func (q *CircularQueue) Full() bool {
	return !q.Empty() && (q.tail+1)%q.length == q.head
}
