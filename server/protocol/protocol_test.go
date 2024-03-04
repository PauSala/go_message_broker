package protocol_test

import (
	"encoding/binary"
	"server/protocol"
	"testing"
)

func TestProtocolReaderShouldParseOperation(t *testing.T) {
	var data [1296]byte = [1296]byte{}
	data[0] |= (0 << 6)
	op, _ := protocol.OpCodeReader(data[0])
	if op != protocol.Pub {
		t.Fatalf("Expected PUB, got %v", op)
	}
	data[0] = 0
	data[0] |= (1 << 6)
	op, _ = protocol.OpCodeReader(data[0])
	if op != protocol.Sub {
		t.Fatalf("Expected SUB, got %v", op)
	}
	data[0] = 0
	data[0] |= (2 << 6)
	op, _ = protocol.OpCodeReader(data[0])
	if op != protocol.Set {
		t.Fatalf("Expected SET, got %v", op)
	}
}

func TestProtocolReaderShouldParsePubMessages(t *testing.T) {
	var data [1296]byte = [1296]byte{}
	data[0] |= (0 << 6)
	s := "A_TOPIC"
	copy(data[1:], []byte(s))
	// Convert the integer to bytes and copy them into data[259:267]
	message := "SOME_MESSAGE"
	l := len(message) * 8
	binary.LittleEndian.PutUint32(data[258:266], uint32(l))
	copy(data[266:266+l], []byte(message))
	//os.WriteFile("pub_message.bin", data[:], 0644)
	res, _ := protocol.ParsePubMessage(data)
	if res.Topic != "A_TOPIC" {
		t.Fatal("Topic not set correctly")
	}
	if res.Data != "SOME_MESSAGE" {
		t.Log(res.Data)
		t.Fatal("Data not set correctly")
	}
}

func TestProtocolReaderShouldParseSetMessges(t *testing.T) {
	var data [1296]byte = [1296]byte{}
	data[0] |= (2 << 6)
	s := "A_TOPIC"
	copy(data[1:], []byte(s))
	//os.WriteFile("set.bin", data[:], 0644)
	res, _ := protocol.ParseSetMessage(data)
	if res.Topic != "A_TOPIC" {
		t.Fatal("Topic not set correctly")
	}
}

func TestProtocolReaderShouldParseSubMessages(t *testing.T) {
	var data [1296]byte = [1296]byte{}
	data[0] |= (1 << 6)
	s := "550e8400-e29b-41d4-a716-446655440000"
	copy(data[1:], []byte(s))
	s = "A_TOPIC"
	copy(data[37:], []byte(s))
	//os.WriteFile("sub.bin", data[:], 0644)
	res, _ := protocol.ParseSubMessage(data)
	if res.SubId != "550e8400-e29b-41d4-a716-446655440000" {
		t.Log(res.SubId)
		t.Fatal("SubId not set correctly")
	}
	if res.Topic != "A_TOPIC" {
		t.Fatal("Topic not set correctly")
	}
	if len(res.SubId) != 36 {
		t.Log(len(res.SubId))
		t.Fatal("Wrong len")
	}
}
