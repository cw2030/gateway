package main

import (
	"fmt"
	"gateway/appcodec"
	"gateway/gw"
	"github.com/cihub/seelog"
)

var (
	addr    = ":3100"
	limit   = 10000
	network = "tcp4"
)

func init() {
	var err error
	gw.Logger, err = seelog.LoggerFromConfigAsFile("seelog.xml")
	if err != nil {
		fmt.Println("Init Log Error:", err)
	} else {
		gw.Logger.Info("Init seelog success.")

	}
}

func main() {
	defer seelog.Flush()
	/*srv := tp.NewPeer(tp.PeerConfig{
		CountTime:  true,
		ListenPort: 9090,
	})*/

	//srv.RouteCall(new(math))
	//srv.ListenAndServe()

	appConf := gw.ServerConf{Network: "tcp4", ServerAddr: ":7722"}
	gw.Logger.Info(appConf)
	app := gw.NewAppMgt(appConf)
	//add Message codec
	//codec := appcodec.SimpleMessageCodec{}
	codec := appcodec.StringMessageCodec{}
	//message handler
	handler := appcodec.AppHandler{}
	app.MessageCodec(codec)
	app.Handler(handler)

	app.Listen()

}
