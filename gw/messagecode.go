package gw

import "net"

type Message interface {
	Encode(connector *Connector) []byte
	Decode([]byte) interface{}
	ToString() string
}
type MessageCodec interface {
	Encode(message Message, connector *Connector) []byte
	Decode(conn net.Conn, connector *Connector) (Message, error)
}
