package main

import (
	"log"
	"os"

	"github.com/bitly/go-nsq"
)

func main() {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer("kefu-dev.hotpu.cn:4150", config)
	if err != nil {
		log.Printf("err:%s\n", err.Error())
		os.Exit(0)
	}
	if err := producer.Publish("connector.1", []byte("hello world")); err != nil {
		log.Printf("err:%s\n", err.Error())
	}
}
