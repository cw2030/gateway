package appcodec

import (
	"errors"
	"gateway/gw"
	"io"
	"net"
)

type StringMessage struct {
	Header *Header
	Body   *Body
}

func (m *StringMessage) Encode() []byte {
	bs := m.Body.ToBytes()
	m.Header.Length = uint16(len(bs))
	total := make([]byte, 12+m.Header.Length)
	copy(total, m.Header.toBytes())
	copy(total[12:], bs)
	return total
}

func (m *StringMessage) Decode([]byte) interface{} {
	panic("implement me")
}

func (m *StringMessage) ToString() string {
	return m.Header.ToString() + m.Body.ToString()
}

type StringMessageCodec struct {
}

func (m StringMessageCodec) Encode(message gw.Message) []byte {
	return message.Encode()
}

func (m StringMessageCodec) Decode(conn net.Conn) (gw.Message, error) {
	headerBytes := make([]byte, headerLength)
	_, err := io.ReadFull(conn, headerBytes)
	if err != nil {
		return nil, errors.New("read Header message fail:" + err.Error())
	}
	h := &Header{}
	h.bytesTo(headerBytes)

	bodyBytes := make([]byte, h.Length)
	io.ReadFull(conn, bodyBytes)
	b := &Body{}
	b.bytesTo(bodyBytes)

	stringMsg := &StringMessage{}
	stringMsg.Header = h
	stringMsg.Body = b
	return stringMsg, nil
}
