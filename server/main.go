package main

import (
	"fmt"
	"server/broker"
	m "server/message"
	"server/protocol"
	"time"
)

func main() {
	s := &broker.Broker{
		Addr: ":3000",
		Dispatcher: *broker.NewDispatcher(
			"DispatcherName",
			broker.SetMaxWorkers(10),
		),
		Parser:    protocol.MessageParser,
		QueryChan: make(chan protocol.Message, 10),
		Queues:    make(map[string]*m.MessageQueue),
	}
	ticker := time.Tick(5 * time.Second)

	go func() {
		for range ticker {
			for _, q := range s.Queues {
				q.PullC <- struct{}{}
			}
		}
	}()

	if err := s.Listen(); err != nil {
		fmt.Println("Server error:", err)
	}
}
