package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	// to produce messages
	topic := "demo"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:29092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}
	defer conn.Close()

	go func() {
		for {
			s := <-signalChanel
			switch s {
			// kill -SIGHUP XXXX [XXXX - PID for your program]
			case syscall.SIGHUP:
				fmt.Println("Signal hang up triggered.")
				if err := conn.Close(); err != nil {
					log.Fatal("failed to close writer:", err)
				}

				// kill -SIGINT XXXX or Ctrl+c  [XXXX - PID for your program]
			case syscall.SIGINT:
				fmt.Println("Signal interrupt triggered.")
				if err := conn.Close(); err != nil {
					log.Fatal("failed to close writer:", err)
				}

				// kill -SIGTERM XXXX [XXXX - PID for your program]
			case syscall.SIGTERM:
				fmt.Println("Signal terminte triggered.")
				if err := conn.Close(); err != nil {
					log.Fatal("failed to close writer:", err)
				}

				// kill -SIGQUIT XXXX [XXXX - PID for your program]
			case syscall.SIGQUIT:
				fmt.Println("Signal quit triggered.")
				if err := conn.Close(); err != nil {
					log.Fatal("failed to close writer:", err)
				}

			default:
				fmt.Println("Unknown signal.")
			}
		}
	}()
	// conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	for {
		_, err = conn.WriteMessages(
			kafka.Message{Value: []byte("one!")},
			kafka.Message{Value: []byte("two!")},
			kafka.Message{Value: []byte("three!")},
		)
		if err != nil {
			log.Fatal("failed to write messages:", err)
		}

		time.Sleep(time.Millisecond)
	}
}
