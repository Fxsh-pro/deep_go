package main

import "errors"

type CircularQueuePointers struct {
	length   int
	capacity int
	head     *QueueElement
	tail     *QueueElement
}

type QueueElement struct {
	value int
	next  *QueueElement
	prev  *QueueElement
}

func NewCircularQueueP(size int) (CircularQueuePointers, error) {
	if size <= 0 {
		return CircularQueuePointers{}, errors.New("size must be greater than 0")
	}
	return CircularQueuePointers{
		capacity: size,
		length:   0,
		head:     nil,
		tail:     nil,
	}, nil
}

func (q *CircularQueuePointers) Push(value int) bool {
	if q.Full() {
		return false
	}
	if q.head == nil {
		q.head = &QueueElement{value: value, prev: nil}
		q.tail = q.head
	} else {
		q.tail.next = &QueueElement{value: value, prev: nil}
		q.tail = q.tail.next
	}
	q.length++
	return true
}

func (q *CircularQueuePointers) Pop() bool {
	if q.Empty() {
		return false
	}
	q.head = q.head.next
	q.length--

	if q.head == nil {
		q.tail = nil
	}

	return true
}

func (q *CircularQueuePointers) Front() int {
	if q.head != nil {
		return q.head.value
	}
	return -1
}

func (q *CircularQueuePointers) Back() int {
	if q.tail != nil {
		return q.tail.value
	}
	return -1
}

func (q *CircularQueuePointers) Empty() bool {
	return q.length == 0
}

func (q *CircularQueuePointers) Full() bool {
	return q.length == q.capacity
}
