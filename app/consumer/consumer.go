package consumer

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/wtifs/room-booking/app/service/log"
	"github.com/wtifs/room-booking/app/service/recovery"
	"github.com/wvanbergen/kafka/consumergroup"
	"time"
)

const (
	BOOKING_STATUS_CG = "momo-q1c28"
	BOOKING_STATUS_TOPIC = ""
)

var (
	zks = []string{
		"10.30.102.60:9092",
	}
)

type roomBookingKafkaConsumer struct {
	consumerGroup *consumergroup.ConsumerGroup
}

func (consumer *roomBookingKafkaConsumer) init() error {
	cfg := consumergroup.NewConfig()
	cfg.Offsets.Initial = sarama.OffsetNewest
	cfg.Offsets.ProcessingTimeout = 10 * time.Second

	var err error
	consumer.consumerGroup, err = consumergroup.JoinConsumerGroup(BOOKING_STATUS_CG, []string{BOOKING_STATUS_TOPIC}, zks, cfg)
	if err != nil {
		return err
	} else {
		log.Info("%s: join consumer group %s successfully", BOOKING_STATUS_TOPIC, BOOKING_STATUS_CG)
	}
	go func() {
		for err := range consumer.consumerGroup.Errors() {
			log.Err("consumer group error: %s", err.Error())
		}
	}()
	return nil
}


//for循环消费
func (consumer *roomBookingKafkaConsumer) consume(ctx context.Context) {
ConsumeMessage:
	for {
		select {
		case <-ctx.Done():
			if err := consumer.consumerGroup.Close(); err != nil {
				log.Err("marketing_user_activation: error closing the consumer: %s", err.Error())
			}
			break ConsumeMessage
		case msg := <-consumer.consumerGroup.Messages():

			//processMsg(ctx, msg.Value)

			//commit after process, confirm at least once
			err := consumer.consumerGroup.CommitUpto(msg)
			if err != nil {
				log.Err("%s: error committing: %s", BOOKING_STATUS_TOPIC, err.Error())
			}
		}
	}
}

//消费activation-affiliate topic, 入库
func RunConsumer() {
	defer recovery.Recovery("run consumer")

	kafkaConsumer := &roomBookingKafkaConsumer{}

	ctx := context.Background()
	err := kafkaConsumer.init()
	if err != nil {
		log.Err("consumer init error: %s", err.Error())
		return
	}
	kafkaConsumer.consume(ctx)
}