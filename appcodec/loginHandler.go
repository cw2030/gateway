package appcodec

import (
	"gateway/gw"
)

/*
处理客户端登录（短信验证、三方登录以及其它方式登录）
*/
func processLogin(message gw.Message, connector *gw.Connector) gw.Message {
	stringMsg := message.(*StringMessage)
	b := stringMsg.Body

	ns := NewEmptyMsg()
	resp := ns.(*StringMessage)
	hr := resp.Header
	br := resp.Body

	hr.ReqType = Response
	hr.MsgType = Msg_Type_Login
	hr.EncryptType = connector.Conf.Encrypt

	br.BType = b.BType
	br.Content = "Login success"
	return resp
}
