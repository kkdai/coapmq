package main

import (
	"log"

	. "github.com/kkdai/coapmq"
)

func main() {
	log.Println("Server start....")
	serv := NewBroker(1024)
	serv.ListenAndServe(":5683")
}
