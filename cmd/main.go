package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
)

var machines = []string { "10.30.102.147", "10.30.102.148", "10.30.102.149" }
type BucketInfo struct {
	Num int
	Lock *sync.Mutex
}

func  main() {
	go kafka()
	getRoomInfo("1")
	roomNum := BucketInfo{
		Num: 0,
	}
	bookNum := BucketInfo{
		Num: 0,
	}
	fmt.Println(roomNum, bookNum)


}


func getRoomInfo(roomId string) string {
	values := map[string]string{"user_token": "6001ff8445378791b8d8f1524660c983f0dc70c8", "meeting_room_id": roomId}
	jsonData, err := json.Marshal(values)
	if err != nil {
		log.Println(err)
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/api/info", machines[rand.Intn(3)]), "application/json", bytes.NewBuffer(jsonData))
	//resp, err := http.Post("http://127.0.0.1:8082", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(content))
	return string(content)
}

func book(roomId string) int {
	values := map[string]string{"user_token": "6001ff8445378791b8d8f1524660c983f0dc70c8", "meeting_room_id": roomId}
	jsonData, err := json.Marshal(values)
	if err != nil {
		log.Println(err)
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/api/book", machines[rand.Intn(3)]), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return 400
	}
	defer resp.Body.Close()
	code := resp.StatusCode
	return code
}

var (
	brokerList        = kingpin.Flag("brokerList", "List of brokers to connect").Default("10.19.1.2:9092").Strings()
	topic             = kingpin.Flag("topic", "Topic name").Default("test").String()
	messageCountStart = kingpin.Flag("messageCountStart", "Message counter start from:").Int()
)

func kafka() {
	kingpin.Parse()
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := *brokerList
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := master.Close(); err != nil {
			log.Panic(err)
		}
	}()
	consumer, err := master.ConsumePartition(*topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Panic(err)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				log.Println(err)
			case msg := <-consumer.Messages():
				*messageCountStart++
				// log.Println("Received messages", string(msg.Key), string(msg.Value))
				log.Println("Received messages", string(msg.Key), string(msg.Value))
				// 处理 kafka 消息
			case <-signals:
				log.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	log.Println("Processed", *messageCountStart, "messages")
}

func dealKafka(value []byte) {
	log.Println(string(value))

}