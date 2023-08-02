package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"log"
	"os"
	"time"
)

type Message struct {
	Content []byte `json:"content"`
}

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Drain()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	var path string
	fmt.Printf("Enter the path to file: ")
	fmt.Scan(&path)
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	message := Message{
		Content: b,
	}

	_, err = js.PublishAsync("subject", message.Content)
	if err != nil {
		log.Println(err)
	}

	select {
	case <-js.PublishAsyncComplete():
		fmt.Println("published complete")
	case <-time.After(time.Second):
		log.Fatal("publish took too long")
	}
}
