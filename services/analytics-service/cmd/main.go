package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
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
		log.Printf("Received message on subject '%s': %s", msg.Subject, string(msg.Data))
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
