package datastructures

import (
	"fmt"
)

type QNode[T comparable] struct {
	Value T
	Next  *QNode[T]
}

type Queue[T comparable] struct {
	Cap  int
	Len  int
	head *QNode[T]
	tail *QNode[T]
}

func NewQueue[T comparable](cap int) *Queue[T] {
	return &Queue[T]{Cap: cap}
}

func (q *Queue[T]) Enqueue(value T) error {
	if q.Len == q.Cap {
		return fmt.Errorf("Queue is full {%v} {%v}", q.Len, q.Cap)
	}
	q.Len += 1
	node := &QNode[T]{Value: value}
	if q.tail == nil {
		q.head = node
		q.tail = node
		return nil
	}

	q.tail.Next = node
	q.tail = node
	return nil
}

func (q *Queue[T]) Dequeue() (T, error) {
	if q.head == nil {
		return Zero[T](), fmt.Errorf("Queue is empty")
	}
	q.Len -= 1
	res := q.head
	q.head = q.head.Next
	if q.head == nil {
		q.tail = nil
	}
	res.Next = nil
	return res.Value, nil
}

func (q *Queue[T]) Peek() T {
	return q.head.Value
}
