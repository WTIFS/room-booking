package main

import (
	"fmt"
	"github.com/wtifs/room-booking/app/consumer"
	"github.com/wtifs/room-booking/app/service/booking"
)

func main() {
	go consumer.RunBookingStatusConsumer()

	booking.Book()
	go dealKafkaChan()
}


func dealKafkaChan() {
	for id := range consumer.KafkaChan {
		fmt.Println(id)
		// 具体的处理逻辑
	}
}