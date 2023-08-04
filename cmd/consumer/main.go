package main

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"log"
	"time"
)

func main() {
	err := readFromNats()
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Println(string(data))

	//msg.Ack()
}

func readFromNats() error {
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Drain()
	streamName := "my_stream"

	js, err := jetstream.New(nc)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	defer js.DeleteStream(ctx, streamName)
	stream, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{"subject"},
	})
	if err != nil {
		return err
	}
	cons, _ := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		InactiveThreshold: 5 * time.Minute,
		FilterSubject:     "subject",
	})
	defer stream.DeleteConsumer(ctx, cons.CachedInfo().Name)
	fetchResult, _ := cons.Fetch(40)
	for msg := range fetchResult.Messages() {
		fmt.Printf("received %s\n", string(msg.Data()))
	}
	return nil
}
