package main

import "net"

type Message interface {
	Encode() []byte
	Decode([]byte) interface{}
	ToString() string
}
type MessageCodec interface {
	Encode(message Message) []byte
	Decode(conn net.Conn) (Message, error)
}
