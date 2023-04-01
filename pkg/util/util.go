// Package util represent few common functions
package util

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// GenerateUUID generate v4 UUID
func GenerateUUID() string {
	u := uuid.New()
	return u.String()
}

// RetryWithExponentialBackoff retries a function with exponential backoff in case of errors
func RetryWithExponentialBackoff(fn func() error, maxRetries int, initialBackoffSeconds int) error {
	backoff := time.Duration(initialBackoffSeconds) * time.Second
	subsequentBackoff := 2

	for i := 0; ; i++ {
		fmt.Println("Retrying....")
		err := fn()
		if err == nil {
			return nil
		}

		if i == maxRetries {
			return fmt.Errorf("maximum number of retries exceeded: %w", err)
		}

		//log.Printf("Error occurred: %v. Retrying in %v", err, backoff)
		time.Sleep(backoff)

		backoff *= time.Duration(subsequentBackoff)
	}
}
