package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

func main() {
	signalChanel := make(chan os.Signal, 1)
	signal.Notify(signalChanel,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	// make a new reader that consumes from topic-A, partition 0, at offset 42
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:29092"},
		Topic:     "demo",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})
	// r.SetOffset(42)

	go func() {
		for {
			s := <-signalChanel
			switch s {
			// kill -SIGHUP XXXX [XXXX - PID for your program]
			case syscall.SIGHUP:
				fmt.Println("Signal hang up triggered.")
				if err := r.Close(); err != nil {
					log.Fatal("failed to close reader:", err)
				}

				// kill -SIGINT XXXX or Ctrl+c  [XXXX - PID for your program]
			case syscall.SIGINT:
				fmt.Println("Signal interrupt triggered.")
				if err := r.Close(); err != nil {
					log.Fatal("failed to close reader:", err)
				}

				// kill -SIGTERM XXXX [XXXX - PID for your program]
			case syscall.SIGTERM:
				fmt.Println("Signal terminte triggered.")
				if err := r.Close(); err != nil {
					log.Fatal("failed to close reader:", err)
				}

				// kill -SIGQUIT XXXX [XXXX - PID for your program]
			case syscall.SIGQUIT:
				fmt.Println("Signal quit triggered.")
				if err := r.Close(); err != nil {
					log.Fatal("failed to close reader:", err)
				}

			default:
				fmt.Println("Unknown signal.")
			}
		}
	}()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Println("error reading message", err)
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
}
