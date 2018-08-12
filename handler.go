package main

type Handler interface {
	HandleFunc(message Message) error
}
