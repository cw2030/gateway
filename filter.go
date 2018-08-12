package main

type Filter interface {
	filter(message Message) error
}
