package main

import (
	"github.com/wtifs/room-booking/app/consumer"
)

func main() {

	c := make(chan bool)

	go consumer.RunConsumer()

	c <- true
}