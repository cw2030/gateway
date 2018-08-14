package main

import "fmt"

type person interface {
	Age() int
}

type Body struct {
}

func (Body) Age() int {
	return 19
}

func main() {
	var p person
	p = Body{}
	fmt.Println(p.Age())
}
