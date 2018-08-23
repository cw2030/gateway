package appcodec

import (
	"gateway/gw"
	"io"
	"net"
)

type StringMessage struct {
	Header *Header
	Body   *Body
}

/*
心跳、握手类型的数据不加密处理，其它均需要加密处理
*/
func (m *StringMessage) Encode(connector *gw.Connector) []byte {
	if m.Header.MsgType == Msg_type_Heartbeat {
		m.Header.Length = 0
		return m.Header.toBytes()
	}

	var bs []byte
	if m.Header.MsgType == Msg_Type_Handshake {
		bs = m.Body.ToBytes()
	} else {
		switch m.Header.EncryptType {
		case gw.Encrypt_AES:
			bs = EncryptAES(connector.Key, m.Body.ToString())
		case gw.Encrypt_RSA2048:
			bs = EncryptRSA2048(m.Body.ToBytes())
		default:
			bs = m.Body.ToBytes()
		}
	}

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

func (m StringMessageCodec) Encode(message gw.Message, connector *gw.Connector) []byte {
	return message.Encode(connector)
}

/*
心跳、握手类型的数据不加密处理，其它均需要加密处理
*/
func (m StringMessageCodec) Decode(conn net.Conn, connector *gw.Connector) (gw.Message, error) {
	defer func() {
		if err := recover(); err != nil {
			gw.Logger.Error(err)
		}
	}()
	headerBytes := make([]byte, headerLength)
	_, err := io.ReadFull(conn, headerBytes)
	if err != nil {
		return nil, err
	}
	h := &Header{}
	h.bytesTo(headerBytes)

	b := &Body{}
	if h.Length > 0 {
		bodyBytes := make([]byte, h.Length)
		_, err = io.ReadFull(conn, bodyBytes)
		if err != nil {
			return nil, err
		}
		if h.ReqType == Response && h.MsgType == Msg_Type_Handshake {
			b.BytesTo(bodyBytes)
		} else {
			switch h.EncryptType {
			case gw.Encrypt_AES:
				b.BytesTo(DecryptAES(connector.Key, bodyBytes))
			case gw.Encrypt_RSA2048:
				b.BytesTo(DecryptRSA2048(bodyBytes))
			default:
				b.BytesTo(bodyBytes)
			}
		}

	}

	stringMsg := &StringMessage{}
	stringMsg.Header = h
	stringMsg.Body = b
	return stringMsg, nil
}
