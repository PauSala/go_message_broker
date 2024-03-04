package messages

import (
	ds "server/data_structures"
)

type MessageQueue struct {
	Id    string
	Queue ds.Queue[string]
	Cap   int
}

func NewMessageQueue(id string, cap int) *MessageQueue {
	return &MessageQueue{
		Id:    id,
		Queue: ds.Queue[string]{Cap: cap, Len: 0},
		Cap:   cap,
	}
}

func (q *MessageQueue) Push(m string) error {
	return q.Queue.Enqueue(m)
}

func (q *MessageQueue) Pull() (string, error) {
	return q.Queue.Dequeue()
}
