package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kubemq-io/kubemq-go"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	eventsStoreClient, err := kubemq.NewEventsStoreClient(ctx,
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId("go-sdk-cookbook-pubsub-events-store-load-balance"),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := eventsStoreClient.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	channel := "events-store.single"

	time.Sleep(100 * time.Millisecond)
	var gracefulShutdown = make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)
	signal.Notify(gracefulShutdown, syscall.SIGQUIT)
	streamSender, err := eventsStoreClient.Stream(ctx, func(result *kubemq.EventStoreResult, err error) {
		if err != nil {
			log.Fatal(err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	counter := 0

	for {
		select {
		case <-time.After(time.Second):
			log.Println("Send")
			counter++
			// err := streamSender(kubemq.NewEventStore().
			// 	SetId("some-id").
			// 	SetChannel(channel).
			// 	SetMetadata("some-metadata").
			// 	SetBody([]byte(fmt.Sprintf("hello kubemq - sending event stream %d", counter))))
			// if err != nil {
			// 	log.Fatal(err)
			// }
			err := streamSender(kubemq.NewEvent().
				SetId("some-id").
				SetChannel(channel).
				SetMetadata("some-metadata").
				AddTag("key", "value").
				SetBody([]byte(fmt.Sprintf("hello kubemq - sending event stream %d", counter))))

			if err != nil {
				log.Fatal(err)
			}
		case <-gracefulShutdown:
			return
		}
	}
}
