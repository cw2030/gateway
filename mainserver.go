package main

import (
	"fmt"
	"github.com/henrylee2cn/teleport"
	"time"
)

var (
	addr    = ":3100"
	limit   = 10000
	network = "tcp4"
)

func main() {
	srv := tp.NewPeer(tp.PeerConfig{
		CountTime:  true,
		ListenPort: 9090,
	})
	srv.RouteCall(new(math))
	srv.ListenAndServe()

	//select {
	//case <-time.After(time.Second * 200):
	//	fmt.Println("end")
	//}

}

type math struct {
	tp.CallCtx
}

func (m *math) Add(args *[]int) (int, *tp.Rerror) {
	if m.Query().Get("push_status") == "yes" {
		m.Session().Push(
			"/push/status",
			fmt.Sprintf("%d numbers are being added...", len(*args)),
		)
		time.Sleep(time.Millisecond * 10)
	}
	var r int
	for _, a := range *args {
		r += a
	}
	return r, nil
}
