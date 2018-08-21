package appcodec

import "gateway/gw"

/*
处理客户端登录（短信验证、三方登录以及其它方式登录）
*/
func processLogin(message gw.Message) gw.Message {
	stringMsg := message.(*StringMessage)
	h := stringMsg.Header
	b := stringMsg.Body
	encrypt := h.EncryptType

}
