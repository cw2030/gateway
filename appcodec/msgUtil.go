package appcodec

import (
	"encoding/hex"
	"gateway/gw"
	"strconv"
	"time"
)

func NewEmptyMsg() gw.Message {
	return &StringMessage{Header: NewHeader(), Body: &Body{}}
}

func NewHandShakeReqMsg() gw.Message {
	msg := &StringMessage{}
	h := NewHeader()
	h.ReqType = Request
	h.MsgType = Msg_Type_Handshake

	body := &Body{}
	msg.Header = h
	msg.Body = body
	return msg
}

func NewHandShakeRespMsg(connector *gw.Connector) gw.Message {
	conf := connector.Conf
	msg := &StringMessage{}
	b := &Body{}
	h := NewHeader()
	h.ReqType = Response
	h.MsgType = Msg_Type_Handshake
	switch conf.Encrypt {
	case gw.Encrypt_None:
		h.EncryptType = gw.Encrypt_None
	case gw.Encrypt_AES:
		h.EncryptType = gw.Encrypt_AES
		b.Content = hex.EncodeToString(connector.Key)
	case gw.Encrypt_RSA1024:
		h.EncryptType = gw.Encrypt_RSA1024
	case gw.Encrypt_RSA2048:
		h.EncryptType = gw.Encrypt_RSA2048
	case gw.Encrypt_SM2:
		h.EncryptType = gw.Encrypt_SM2
	case gw.Encrypt_SM4:
		h.EncryptType = gw.Encrypt_SM4
	default:
		h.EncryptType = gw.Encrypt_None
	}
	msg.Header = h
	msg.Body = b
	return msg
}

func NewHeatbeatReqMsg() gw.Message {
	msg := &StringMessage{}
	b := &Body{}
	h := NewHeader()
	h.ReqType = Request
	h.MsgType = Msg_type_Heartbeat

	b.Content = strconv.FormatInt(time.Now().Unix(), 10)

	msg.Header = h
	msg.Body = b
	return msg
}
func NewHeatbeatRespMsg() gw.Message {
	msg := &StringMessage{}
	b := &Body{}
	h := NewHeader()
	h.ReqType = Response
	h.MsgType = Msg_type_Heartbeat

	b.Content = strconv.FormatInt(time.Now().Unix(), 10)

	msg.Header = h
	msg.Body = b
	return msg
}

func NewLoginReqMsg() gw.Message {
	msg := &StringMessage{}

	return msg
}

func NewLoginRespMsg() gw.Message {
	msg := &StringMessage{}

	return msg
}
