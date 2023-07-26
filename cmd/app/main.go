package main

import (
	"github.com/nats-io/nats.go"
	"log"
)

func main() {
	// Подключение к серверу Nats
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatal(err)
	}
	sub, err := js.SubscribeSync("subject")
	if err != nil {
		log.Fatal(err, "smth")
	}
	for {
		msg, err := sub.NextMsg(nats.DefaultTimeout)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Получено сообщение: %s", string(msg.Data))
	}
}
