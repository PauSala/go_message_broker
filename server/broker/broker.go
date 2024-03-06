package broker

import (
	"fmt"
	"net"
	messages "server/message"
	"server/protocol"
	"strings"
)

type Broker struct {
	Addr       string
	Dispatcher Dispatcher
	Parser     protocol.Parser
	Queues     map[string]*messages.MessageQueue
	QueryChan  chan protocol.Message
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
	q := messages.NewMessageQueue(name, 10)
	q.Listen()
	s.Queues[name] = q
}

func (s *Broker) MessageListener() {
	for m := range s.QueryChan {
		kind := m.Kind()
		switch kind {
		case protocol.Pub:
			fmt.Println("CHANNEL: Received a Pub message")
			fmt.Println(m.(*protocol.PubMessage).Topic)
			fmt.Println(m.(*protocol.PubMessage).Data)
			queue := s.Queues[m.(*protocol.PubMessage).Topic]
			if queue == nil {
				fmt.Println("CHANNEL: Queue not found")
				break
			}
			queue.PushC <- m.(*protocol.PubMessage).Data
		case protocol.Sub:
			fmt.Println("CHANNEL: Received a Sub message")
			queue := s.Queues[m.(*protocol.SubMessage).Topic]
			if queue == nil {
				fmt.Println("CHANNEL: Queue not found")
				break
			}
			ip := net.ParseIP(m.(*protocol.SubMessage).SubIp)
			if queue.SubC != nil {
				queue.SubC <- &messages.Subscriber{
					Id: m.(*protocol.SubMessage).SubId,
					Ip: net.TCPAddr{IP: ip, Port: 3002}}
			}
		case protocol.Set:
			fmt.Println("CHANNEL: Received a Set message")
			queue := s.Queues[m.(*protocol.SetMessage).Topic]
			if queue == nil {
				s.AddQueue(m.(*protocol.SetMessage).Topic)
				break
			}
			fmt.Println("CHANNEL: Queue already exists")
		case protocol.Pull:
			fmt.Println("CHANNEL: Recieved a Pull message")
			queue := s.Queues[m.(*protocol.PullMessage).Topic]
			if queue == nil {
				fmt.Println("CHANNEL: Queue not found")
				break
			}
			queue.PullC <- struct{}{}
		}
	}
}

func (s *Broker) Listen() error {
	listener, err := net.Listen("tcp", s.Addr)
	go s.MessageListener()
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
		m := m.(*protocol.SubMessage)
		m.SubIp = strings.Split(remoteAddr, ":")[0]
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
