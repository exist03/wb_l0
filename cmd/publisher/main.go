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

	path := ""
	fmt.Printf("Enter the path to file: ")
	fmt.Scan(&path)
	b, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	message := Message{
		Content: b,
	}

	// Отправка сообщения в JetStream
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

//func ExampleJetStream() {
//	streamName := "EVENTS"
//	nc, err := nats.Connect(nats.DefaultURL)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer nc.Drain()
//
//	newJS, _ := jetstream.New(nc)
//
//	// The new API uses `context.Context` for cancellation and timeouts when
//	// managing streams and consumers.
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	// Creating a stream is done using the `CreateStream` method.
//	// It works similarly to the legacy `AddStream` method, except
//	// instead of returning `StreamInfo`, it returns a `Stream` handle,
//	// which can be used to manage the stream.
//	// Instead of creating a new stream, let's look up the existing `EVENTS` stream.
//	stream, _ := newJS.Stream(ctx, streamName)
//
//	// The new API differs from the legacy API in that it does not
//	// auto-create consumers. Instead, consumers must be created or retrieved
//	// explicitly. This allows for more control over the consumer lifecycle,
//	// while also getting rid of the hidden logic of the `Subscribe()` methods.
//	// In order to create a consumer, use the `AddConsumer` method.
//	// This method works similarly to the legacy `AddConsumer` method,
//	// except it returns a `Consumer` handle, which can be used to manage
//	// the consumer. Notice that since we are using pull consumers, we
//	// do not need to provide a `DeliverSubject`.
//	// In order to create a short-lived, ephemeral consumer, we will set the
//	// `InactivityThreshold` to a low value and not provide a consumer name.
//	cons, _ := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
//		InactiveThreshold: 10 * time.Millisecond,
//	})
//	fmt.Println("Created consumer", cons.CachedInfo().Name)
//
//	// ### Continuous message retrieval with `Consume()`
//	// In order to continuously receive messages, the `Consume` method
//	// can be used. This method works similarly to the legacy `Subscribe`
//	// method, in that it will asynchronously deliver messages to the
//	// provided `jetstream.MsgHandler` function. However, it does not
//	// create a consumer, instead it will use the consumer created
//	// previously.
//	fmt.Println("# Consume messages using Consume()")
//	consumeContext, _ := cons.Consume(func(msg jetstream.Msg) {
//		fmt.Printf("received %q\n", msg.Subject())
//		msg.Ack()
//	})
//	time.Sleep(100 * time.Millisecond)
//
//	// `Consume()` returns a `jetstream.ConsumerContext` which can be used
//	// to stop consuming messages. In contrast to `Unsubscribe()` in the
//	// legacy API, this will not delete the consumer.
//	// Consumer will be automatically deleted by the server when the
//	// `InactivityThreshold` is reached.
//	consumeContext.Stop()
//
//	// Now let's create a new, long-lived, named consumer.
//	// In order to filter messages, we will provide a `FilterSubject`.
//	// This is equivalent to providing a subject to `Subscribe` in the
//	// legacy API.
//	// InactiveThreshold will cause the consumer to be automatically
//	// removed after 10 minutes of inactivity. It can be omitted
//	// for durable consumers.
//	consumerName := "pull-1"
//	cons, _ = stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
//		Name:              consumerName,
//		InactiveThreshold: 10 * time.Minute,
//		FilterSubject:     "events.2",
//	})
//	fmt.Println("Created consumer", cons.CachedInfo().Name)
//
//	// As an alternative to `Consume`, the `Messages()` method can be used
//	// to retrieve messages one-by-one. Note that this method will
//	// still pre-fetch messages, but instead of delivering them to a
//	// handler function, it will return them upon calling `Next`.
//	fmt.Println("# Consume messages using Messages()")
//	it, _ := cons.Messages()
//	msg1, _ := it.Next()
//	fmt.Printf("received %q\n", msg1.Subject())
//
//	// Similarly to `Consume`, `Messages` allows to stop consuming messages
//	// without deleting the consumer.
//	it.Stop()
//
//	// ### Retrieving messages on demand with `Fetch()` and `Next()`
//
//	// Similar to the legacy API, the new API also exposes a `Fetch()`
//	// method for retrieving a specified number of messages on demand.
//	// This method resembles the legacy `FetchBatch` method, in that
//	// it will return a channel on which the messages will be delivered.
//	fmt.Println("# Fetch messages")
//	cons, _ = stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
//		InactiveThreshold: 10 * time.Millisecond,
//	})
//	fetchResult, _ := cons.Fetch(2, jetstream.FetchMaxWait(100*time.Millisecond))
//	for msg := range fetchResult.Messages() {
//		fmt.Printf("received %q\n", msg.Subject())
//		msg.Ack()
//	}
//
//	// Alternatively, the `Next` method can be used to retrieve a single
//	// message. It works like `Fetch(1)`, returning a single message instead
//	// of a channel.
//	fmt.Println("# Get next message")
//	msg1, _ = cons.Next()
//	fmt.Printf("received %q\n", msg1.Subject())
//	msg1.Ack()
//
//	// Streams and consumers can be deleted using the `DeleteStream` and
//	// `DeleteConsumer` methods. Note that deleting a stream will also
//	// delete all consumers on that stream.
//	fmt.Println("# Delete consumer")
//	stream.DeleteConsumer(ctx, cons.CachedInfo().Name)
//	fmt.Println("# Delete stream")
//	newJS.DeleteStream(ctx, streamName)
//}
