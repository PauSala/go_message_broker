package protocol

import (
	"encoding/binary"
	"errors"
	"strings"

	"github.com/google/uuid"
)

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

// Deffinition
/*
Operations:
- Subscribe: subscriberId (272 bits) Topic (2048 bits)

0    4    8        16        24        32
+----+----+---------+---------+---------+
| Op |                                  |
+---------+---------+---------+---------+
|          subscriberID (272)           |
+---------+---------+---------+---------+
|          Topic (2048)                 |
+---------+---------+---------+---------+

- Publish: Topic (1024 bits) (length + body)

0    4    8        16        24        32
+----+----+---------+---------+---------+
| Op |                                  |
+---------+---------+---------+---------+
|          Topic (2048)                 |
+---------+---------+---------+---------+
|          length                       |
+---------+---------+---------+---------+
|                                       |
.           ...  body ...               .
.                                       .
.                                       .
+----------------------------------------


- Set Topic Topic (2048)

0    4    8        16        24        32
+----+----+---------+---------+---------+
| Op |                                  |
+---------+---------+---------+---------+
|          Topic (2048)                 |
+---------+---------+---------+---------+
*/

type MessageKind int

const (
	Pub MessageKind = iota
	Sub
	Set
)

func (mk MessageKind) String() string {
	switch mk {
	case Pub:
		return "Pub"
	case Sub:
		return "Sub"
	case Set:
		return "Set"
	default:
		return "Unknown"
	}
}

type Message interface {
	Kind() MessageKind
}

type PubMessage struct {
	Topic string
	Data  string
}

func (m *PubMessage) Kind() MessageKind {
	return Pub
}

type SubMessage struct {
	SubId string
	Topic string
}

func (m *SubMessage) Kind() MessageKind {
	return Sub
}

type SetMessage struct {
	Topic string
}

func (m *SetMessage) Kind() MessageKind {
	return Set
}

type Parser func(message [1296]byte) (Message, error)

func MessageParser(message [1296]byte) (Message, error) {

	opCode, err := OpCodeReader(message[0])
	if err != nil {
		return &SetMessage{}, err
	}
	switch opCode {
	case Pub:
		return ParsePubMessage(message)
	case Sub:
		return ParseSubMessage(message)
	case Set:
		return ParseSetMessage(message)
	default:
		return &SetMessage{}, errors.New("unknown operation")
	}
}

func ParseSetMessage(message [1296]byte) (*SetMessage, error) {
	s := strings.TrimRight(string(message[1:258]), "\x00")
	return &SetMessage{Topic: s}, nil
}

func ParsePubMessage(Message [1296]byte) (*PubMessage, error) {
	topic := strings.TrimRight(string(Message[1:258]), "\x00")
	len := Message[258:266]
	i := binary.LittleEndian.Uint32(len[:]) / 8
	body := Message[266 : 266+i]
	return &PubMessage{Topic: topic, Data: string(body)}, nil
}

func ParseSubMessage(Message [1296]byte) (*SubMessage, error) {
	subId := string(Message[1:37])
	isValid := IsValidUUID(subId)
	if !isValid {
		return &SubMessage{}, errors.New("invalid UUid")
	}
	topic := strings.TrimRight(string(Message[37:293]), "\x00")
	return &SubMessage{SubId: subId, Topic: topic}, nil
}

func OpCodeReader(c byte) (MessageKind, error) {
	firstTwoBits := c >> 6
	switch firstTwoBits {
	case 0:
		return Pub, nil
	case 1:
		return Sub, nil
	case 2:
		return Set, nil
	default:
		return 0, errors.New("unknown operation")
	}
}
