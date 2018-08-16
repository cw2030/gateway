package main

import (
	"fmt"
	"github.com/cihub/seelog"
)

type person interface {
	Age() int
}

type Body struct {
}

func (Body) Age() int {
	return 19
}

func main() {
	defer seelog.Flush()
	//logger, err := seelog.LoggerFromConfigAsFile("seelog.xml")
	if err != nil {
		fmt.Println(err)
	}
	var p person
	p = Body{}
	fmt.Println(p.Age())
	seelog.Info("abc")
}
