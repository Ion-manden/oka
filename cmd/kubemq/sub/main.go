package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kubemq-io/kubemq-go"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// eventsStoreClient, err := kubemq.NewEventsStoreClient(ctx,
	// 	kubemq.WithAddress("localhost", 50000),
	// 	// kubemq.WithClientId("go-sdk-cookbook-pubsub-events-store-load-balance"),
	// 	kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer func() {
	// 	err := eventsStoreClient.Close()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()
	eventsClient, err := kubemq.NewEventsClient(ctx,
		kubemq.WithAddress("localhost", 50000),
		// kubemq.WithClientId("go-sdk-cookbook-pubsub-events-stream"),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := eventsClient.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	channel := "events-store.single"
	// err = eventsStoreClient.Subscribe(ctx, &kubemq.EventsStoreSubscription{
	// 	Channel: channel,
	// 	// ClientId:         "go-sdk-cookbook-pubsub-events-store-single-subscriber",
	// 	SubscriptionType: kubemq.StartFromFirstEvent(),
	// }, func(msg *kubemq.EventStoreReceive, err error) {
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	} else {
	// 		log.Printf("Receiver - Event Store Received:\nEventID: %s\nChannel: %s\nMetadata: %s\nBody: %s\n", msg.Id, msg.Channel, msg.Metadata, msg.Body)
	// 	}
	// })
	// if err != nil {
	// 	log.Fatal(err)
	// }

	err = eventsClient.Subscribe(ctx, &kubemq.EventsSubscription{
		Channel: channel,
		// ClientId: "go-sdk-cookbook-pubsub-events-stream-subscriber",
	}, func(msg *kubemq.Event, err error) {
		if err != nil {
			log.Fatal(err)
		} else {
			log.Printf("Receiver - Event Received:\nEventID: %s\nChannel: %s\nMetadata: %s\nBody: %s\n", msg.Id, msg.Channel, msg.Metadata, msg.Body)
		}
	})

	var gracefulShutdown = make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)
	signal.Notify(gracefulShutdown, syscall.SIGQUIT)

	for {
		select {
		case <-gracefulShutdown:
			return
		}
	}
}
