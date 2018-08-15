package appcodec

import (
	"gateway/gw"
	"github.com/cihub/seelog"
)

type AppHandler struct {
}

func (AppHandler) HandleFunc(connector *gw.Connector, message gw.Message, err error) {
	sm := message.(*StringMessage)
	seelog.Info("Header:", sm.Header.ToString())
	seelog.Infof("Body:%s", sm.Body.ToString())
	connector.WriteChan <- message.Encode()
}
