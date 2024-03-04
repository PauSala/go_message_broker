package broker

import (
	"fmt"
	"net"
	messages "server/message"
	"server/protocol"
)

type Broker struct {
	Addr         string
	Dispatcher   Dispatcher
	Parser       protocol.Parser
	Queues       map[string]*messages.MessageQueue
	QueryChan    chan protocol.Message
	ResponseChan chan string
}

type Handler func(conn net.Conn) error

type Job struct {
	Broker *Broker
	conn   net.Conn
}

func (s *Job) Process() error {
	return s.Broker.handle(s.conn)
}

func (s *Broker) AddQueue(name string) {
	fmt.Println("Adding queue")
	s.Queues[name] = messages.NewMessageQueue(name, 10)
}

func (s *Broker) MessageListener() {
	for m := range s.QueryChan {
		kind := m.Kind()
		switch kind {
		case protocol.Pub:
			fmt.Println("CHANNEL: Received a Pub message")
			queue := s.Queues[m.(*protocol.PubMessage).Topic]
			if queue == nil {
				fmt.Println("CHANNEL: Queue not found")
				break
			}
			queue.Push(m.(*protocol.PubMessage).Data)
		case protocol.Sub:
			fmt.Println("CHANNEL: Received a Sub message")
			queue := s.Queues[m.(*protocol.SubMessage).Topic]
			if queue == nil {
				fmt.Println("CHANNEL: Queue not found")
				break
			}
			data, err := queue.Pull()
			if err != nil {
				fmt.Println("CHANNEL: Error pulling data from queue: " + err.Error())
				break
			}
			fmt.Println("Data: " + data)
			s.ResponseChan <- data
		case protocol.Set:
			fmt.Println("CHANNEL: Received a Set message")
			fmt.Println(m.(*protocol.SetMessage).Topic)
			queue := s.Queues[m.(*protocol.SetMessage).Topic]
			if queue == nil {
				s.AddQueue(m.(*protocol.SetMessage).Topic)
				break
			}
			fmt.Println("CHANNEL: Queue already exists")
		}
	}
}

func (s *Broker) MessageSender() {
	for msg := range s.ResponseChan {
		fmt.Println("Sending response to client")
		fmt.Println(msg)
	}
}

func (s *Broker) Listen() error {
	listener, err := net.Listen("tcp", s.Addr)
	go s.MessageListener()
	go s.MessageSender()
	s.ResponseChan <- "Some test data"
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		s.Dispatcher.JobQueue <- &Job{conn: conn, Broker: s}
	}
}

func (s *Broker) handle(conn net.Conn) error {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	fmt.Println("Client connected from " + remoteAddr)

	buf := make([]byte, 1296)

	_, err := conn.Read(buf)
	if err != nil {
		return sendError(conn, "Error reading request"+err.Error())
	}

	if len(buf) < 1296 {
		return sendError(conn, "Error: buffer is too small")
	}

	var arr [1296]byte
	copy(arr[:], buf[0:1296])
	m, err := s.Parser(arr)
	if err != nil {
		return sendError(conn, "Error parsing message"+err.Error())
	}
	kind := m.Kind()
	switch kind {
	case protocol.Pub:
		fmt.Println("WORKER: Received a Pub message")
		s.QueryChan <- m
	case protocol.Sub:
		fmt.Println("WORKER: Received a Sub message")
		s.QueryChan <- m
	case protocol.Set:
		fmt.Println("WORKER: Received a Set message")
		s.QueryChan <- m
	}
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 0\r\n" +
		"\r\n"

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Error writing to connection:", err)
		return fmt.Errorf("error writing to connection")
	}
	return nil
}

func sendError(conn net.Conn, e string) error {
	response := "HTTP/1.1 500 Internal Server Error\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 23\r\n" +
		"\r\n" +
		e
	_, err := conn.Write([]byte(response))
	return err
}
