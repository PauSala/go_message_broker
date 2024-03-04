package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"github.com/google/uuid"
)

func listen() error {
	println("TCP server listening on port 3002")
	ln, err := net.Listen("tcp", ":3002")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
			return err
		}
		data := make([]byte, 1024)
		_, err = conn.Read(data)
		if err != nil {
			log.Fatal(err)
			return err
		}
		println(string(data))
	}
}

func main() {
	serverID := uuid.New()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Read the HTTP request body
		data, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// Connect to the TCP server
		conn, err := net.Dial("tcp", "localhost:3000")
		if err != nil {
			http.Error(w, "Failed to connect to TCP server", http.StatusInternalServerError)
			return
		}
		defer conn.Close()

		// Send the HTTP request data to the TCP server
		message := string(data)
		parsed := Build(message, serverID.String())
		if _, err := conn.Write(parsed); err != nil {
			http.Error(w, "Failed to send data to TCP server", http.StatusInternalServerError)
			return
		}

		// Read the response from the TCP server
		response, err := io.ReadAll(conn)
		if err != nil {
			http.Error(w, "Failed to read response from TCP server", http.StatusInternalServerError)
			return
		}

		// Write the TCP server's response as the HTTP response
		w.Write(response)
	})

	go func() {
		if err := listen(); err != nil {
			fmt.Println("Server error:", err)
		}
	}()
	log.Fatal(http.ListenAndServe(":3001", nil))
}
