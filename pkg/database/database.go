package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(databaseURL string) (*pgxpool.Pool, error) {
	var err error
	var pool *pgxpool.Pool

	const maxRetries = 5
	const backoff = 2 * time.Second

	for i := 0; i < maxRetries; i++ {
		// Create a new context for each attempt.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		pool, err = pgxpool.New(ctx, databaseURL)
		if err == nil {
			// If pool is created, try to ping.
			err = pool.Ping(ctx)
			if err == nil {
				// Success!
				cancel()
				log.Println("database connection pool established")
				return pool, nil
			}
			// Ping failed, close the pool before retrying.
			pool.Close()
		}

		// Always cancel the context for the current attempt.
		cancel()

		if i < maxRetries-1 {
			log.Printf("database connection failed (attempt %d/%d), retrying in %v: %v", i+1, maxRetries, backoff, err)
			time.Sleep(backoff)
		}
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
}
