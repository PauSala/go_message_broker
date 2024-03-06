package messages

import (
	"net"
	ds "server/data_structures"
)

type Subscriber struct {
	Id string
	Ip net.TCPAddr
}

type MessageQueue struct {
	Id          string
	Queue       ds.Queue[string]
	Cap         int
	PushC       chan string
	PullC       chan struct{}
	SubC        chan *Subscriber
	Subscribers map[string]*Subscriber
	onPull      onPull
}

type onPull func(msg string, q *MessageQueue)

func NewMessageQueue(id string, cap int) *MessageQueue {
	return &MessageQueue{
		Id:          id,
		Queue:       ds.Queue[string]{Cap: cap, Len: 0},
		Cap:         cap,
		PushC:       make(chan string, 10),
		PullC:       make(chan struct{}, 10),
		SubC:        make(chan *Subscriber, 10),
		Subscribers: make(map[string]*Subscriber),
		onPull:      OnPull,
	}
}

func OnPull(msg string, q *MessageQueue) {
	for _, s := range q.Subscribers {
		conn, err := net.Dial("tcp", s.Ip.String())
		if err != nil {
			println("Failed to connect to TCP server" + err.Error())
			return
		}
		defer conn.Close()
		message := []byte(msg)
		if _, err := conn.Write(message); err != nil {
			println("Failed to send data to TCP server" + err.Error())
			return
		}
	}
}

func (q *MessageQueue) Listen() {
	go func() {
		for {
			select {
			case m := <-q.PushC:
				q.Queue.Enqueue(m)
			case <-q.PullC:
				for msg, err := q.Queue.Dequeue(); err == nil; msg, err = q.Queue.Dequeue() {
					go q.onPull(msg, q)
				}
			case s := <-q.SubC:
				q.Subscribers[s.Id] = s
			}
		}
	}()
}

func (q *MessageQueue) Push(m string) error {
	return q.Queue.Enqueue(m)
}

func (q *MessageQueue) Pull() (string, error) {
	return q.Queue.Dequeue()
}
