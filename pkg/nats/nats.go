package nats

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func ConnectNATS(natsURL string) (*nats.Conn, error) {
	var nc *nats.Conn
	var err error

	const maxRetries = 5
	const backoff = 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		nc, err = nats.Connect(natsURL)
		if err == nil {
			log.Println("connected to nats server")
			return nc, nil
		}

		if i < maxRetries-1 {
			log.Printf("nats connection failed (attempt %d/%d), retrying in %v: %v", i+1, maxRetries, backoff, err)
			time.Sleep(backoff)
		}
	}

	return nil, fmt.Errorf("failed to connect to nats after %d attempts: %w", maxRetries, err)
}
