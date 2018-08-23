package appcodec

import (
	"encoding/json"
	"gateway/gw"
	"net"
)

type AppRespBody struct {
	Status  int    `json:"s"`
	Content string `json:"c"`
	ErrMsg  string `json:"e"`
}

func (ar *AppRespBody) Encode(connector *gw.Connector) []byte {
	bs, err := json.Marshal(ar)
	if err != nil {
		return nil
	}
	return bs
}

func (ar *AppRespBody) Decode([]byte) interface{} {
	panic("implement me")
}

func (ar *AppRespBody) ToString() string {
	bs, err := json.Marshal(ar)
	if err != nil {
		return ""
	}
	return gw.Byte2str(bs)
}

type AppRespMessage struct {
	Header Header
	resp   AppRespBody
}

func (arm *AppRespMessage) Encode(message gw.Message, connector *gw.Connector) []byte {
	if arm.Header.MsgType == Msg_type_Heartbeat {
		arm.Header.Length = 0
		return arm.Header.toBytes()
	}

	var bs []byte
	if arm.Header.MsgType == Msg_Type_Handshake {
		bs = arm.resp.Encode(connector)
	} else {
		switch arm.Header.EncryptType {
		case gw.Encrypt_AES:
			bs = EncryptAES(connector.Key, arm.resp.ToString())
		case gw.Encrypt_RSA2048:
			bs = EncryptRSA2048(arm.resp.Encode(connector))
		default:
			bs = arm.resp.Encode(connector)
		}
	}

	arm.Header.Length = uint16(len(bs))
	total := make([]byte, 12+arm.Header.Length)
	copy(total, arm.Header.toBytes())
	copy(total[12:], bs)
	return total
}

func (arm *AppRespMessage) Decode(conn net.Conn, connector *gw.Connector) (gw.Message, error) {
	panic("implement me")
}
