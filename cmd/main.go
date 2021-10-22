package main

import (
	"github.com/wtifs/room-booking/app/consumer"
	"github.com/wtifs/room-booking/app/service/booking"
)

func main() {
	go consumer.RunBookingStatusConsumer()

	booking.Book()
}
