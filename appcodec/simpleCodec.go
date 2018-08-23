package appcodec

import (
	"errors"
	"fmt"
	"gateway/gw"
	"io"
	"net"
)

type SimpleMessage struct {
	Header *SimpleHeader
	Body   *Body
}

func (m *SimpleMessage) Encode(connector *gw.Connector) []byte {
	bs := m.Body.ToBytes()
	bodyLen := len(bs)
	m.Header.BodyLength = uint16(bodyLen)
	total := make([]byte, 4+bodyLen)

	copy(total, m.Header.toBytes())
	copy(total[4:], bs)
	return total
}

func (m *SimpleMessage) Decode([]byte) interface{} {
	panic("implement me")
}

func (m *SimpleMessage) ToString() string {
	return m.Header.ToString() + m.Body.ToString()
}

type SimpleMessageCodec struct {
}

func (m SimpleMessageCodec) Encode(message gw.Message, connector *gw.Connector) []byte {
	return message.Encode(connector)
}

func (m SimpleMessageCodec) Decode(conn net.Conn, connector *gw.Connector) (gw.Message, error) {
	headerBytes := make([]byte, simpleMessageHeaderLength)
	_, err := io.ReadFull(conn, headerBytes)
	if err != nil {
		return nil, errors.New("read Header message fail:" + err.Error())
	}
	fmt.Println("headerbytes:", headerBytes)
	h := &SimpleHeader{}
	h.bytesTo(headerBytes)

	bodyBytes := make([]byte, h.BodyLength)
	io.ReadFull(conn, bodyBytes)
	b := &Body{}
	b.BytesTo(bodyBytes)

	stringMsg := &SimpleMessage{}
	stringMsg.Header = h
	stringMsg.Body = b
	return stringMsg, nil
}
