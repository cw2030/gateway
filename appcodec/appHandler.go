package appcodec

import (
	"gateway/gw"
	"time"
)

type AppHandler struct {
}

func (AppHandler) HandleFunc(connector *gw.Connector, message gw.Message, err error) {
	sm := message.(*StringMessage)
	gw.Logger.Debugf("Receive Body:%s", sm.ToString())
	switch sm.Header.MsgType {
	case Msg_Type_Handshake:
		connector.WriteChan <- NewHandShakeRespMsg(connector).Encode(connector)
	case Msg_type_Heartbeat:
		connector.LatestActivity = time.Now()
		connector.WriteChan <- NewHeatbeatRespMsg().Encode(connector)
	case Msg_Type_Login:
		connector.WriteChan <- processLogin(sm, connector).Encode(connector)
	case Msg_Type_Busi:
		connector.WriteChan <- proxy(sm, connector).Encode(connector)
	case Msg_Type_Async:
		// Async异步的消息只返回收到请求的应答
		connector.WriteChan <- proxy(sm, connector).Encode(connector)
	default:
		gw.Logger.Warnf("Unknown message Type:%s", sm.ToString())

	}

}
