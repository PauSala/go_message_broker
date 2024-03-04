package main

import (
	"encoding/binary"
	"strings"
)

func Build(query, uuid string) []byte {
	args := strings.Fields(query)
	command := args[0]

	switch command {
	case "PUB":
		msg := strings.Join(args[1:], " ")
		return PubMessage(args[1], msg)
	case "SUB":
		return SubMessage(args[1], uuid)
	case "SET":
		return SetMessage(args[1])
	default:
		return []byte("INVALID COMMAND")
	}
}

func PubMessage(topic, message string) []byte {
	var data [1296]byte = [1296]byte{}
	data[0] |= (0 << 6)
	copy(data[1:], []byte(topic))
	l := len(message) * 8
	binary.LittleEndian.PutUint32(data[258:266], uint32(l))
	copy(data[266:266+l], []byte(message))
	return data[:]
}

func SubMessage(topic, id string) []byte {
	var data [1296]byte = [1296]byte{}
	data[0] |= (1 << 6)
	copy(data[1:], []byte(id))
	copy(data[37:], []byte(topic))
	return data[:]
}

func SetMessage(topic string) []byte {
	var data [1296]byte = [1296]byte{}
	data[0] |= (2 << 6)
	copy(data[1:], []byte(topic))
	return data[:]
}
