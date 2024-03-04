package main

import (
	"fmt"
	"server/broker"
	m "server/message"
	"server/protocol"
)

func main() {
	s := &broker.Broker{
		Addr: ":3000",
		Dispatcher: *broker.NewDispatcher(
			"DispatcherName",
			broker.SetMaxWorkers(10),
		),
		Parser:       protocol.MessageParser,
		QueryChan:    make(chan protocol.Message, 10),
		ResponseChan: make(chan string, 10),
		Queues:       make(map[string]*m.MessageQueue),
	}
	if err := s.Listen(); err != nil {
		fmt.Println("Server error:", err)
	}
}
