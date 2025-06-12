package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	eventsv1 "github.com/nouvadev/veritas/pkg/gen/proto/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

func main() {
	// Connect to NATS
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	defer nc.Close()

	log.Println("Connected to NATS server at", natsURL)

	// Subscribe to the subject
	subject := "veritas.redirect.success"
	sub, err := nc.Subscribe(subject, func(msg *nats.Msg) {
		event := &eventsv1.RedirectEvent{}
		if err := proto.Unmarshal(msg.Data, event); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			return
		}
		log.Printf(
			"Received Event: ShortCode=%s, OriginalURL=%s, UserAgent=%s, IP=%s",
			event.ShortCode,
			event.OriginalUrl,
			event.UserAgent,
			event.IpAddress,
		)
	})
	if err != nil {
		log.Fatalf("Error subscribing to subject '%s': %v", subject, err)
	}
	defer sub.Unsubscribe()

	log.Printf("Subscribed to subject '%s'", subject)

	// Keep the service running
	log.Println("Analytics service is running. Waiting for events...")
	// Wait for a signal to exit
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down analytics service.")
}
