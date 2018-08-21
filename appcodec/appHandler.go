package appcodec

import (
	"gateway/gw"
	"strconv"
	"time"
)

type AppHandler struct {
}

func (AppHandler) HandleFunc(connector *gw.Connector, message gw.Message, err error) {
	sm := message.(*StringMessage)
	gw.Logger.Debugf("Body:%s", sm.ToString())
	switch sm.Header.MsgType {
	case Msg_Type_Handshake:
		connector.WriteChan <- NewHandShakeRespMsg(connector).Encode()
	case Msg_type_Heartbeat:
		ts, err := strconv.ParseInt(sm.Body.Content, 64, 64)
		if err != nil {
			connector.LatestActivity = time.Now()
		} else {
			connector.LatestActivity = time.Unix(ts, 0)
		}
		connector.WriteChan <- NewHeatbeatRespMsg().Encode()
	case Msg_Type_Login:

	case Msg_Type_Busi:
	default:

	}
}
