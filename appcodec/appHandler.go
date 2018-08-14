package appcodec

import (
	"fmt"
	"gateway/gw"
)

type AppHandler struct {
}

func (AppHandler) HandleFunc(connector *gw.Connector, message gw.Message, err error) {
	fmt.Println(message.ToString())
	connector.Conn.Write(message.Encode())
}
