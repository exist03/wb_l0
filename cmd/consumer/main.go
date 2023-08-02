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

func readFromNatsDeprecated() ([]byte, error) {
	nc, _ := nats.Connect(nats.DefaultURL)
	defer nc.Drain()
	streamName := "my_stream"

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{"subject"},
	})
	if err != nil {
		return nil, err
	}
	consumerName := "pull-1"
	cons, _ := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:              consumerName,
		InactiveThreshold: 5 * time.Minute,
		FilterSubject:     "subject",
	})
	fmt.Println("Created consumer", cons.CachedInfo().Name)

	it, _ := cons.Messages()

	msg1, _ := it.Next()
	fmt.Printf("received %q\n", string(msg1.Data()))

	it.Stop()

	fmt.Println("# Delete consumer")
	stream.DeleteConsumer(ctx, cons.CachedInfo().Name)
	fmt.Println("# Delete stream")
	js.DeleteStream(ctx, streamName)
	return msg1.Data(), nil
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
		//consumerName := "pull-1"
	})
	if err != nil {
		return err
	}
	cons, _ := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		//Name:              consumerName,
		InactiveThreshold: 5 * time.Minute,
		FilterSubject:     "subject",
	})
	defer stream.DeleteConsumer(ctx, cons.CachedInfo().Name)
	fetchResult, _ := cons.Fetch(40)
	for msg := range fetchResult.Messages() {
		fmt.Printf("received %s\n", string(msg.Data()))
		//msg.Ack()
	}
	return nil
}
