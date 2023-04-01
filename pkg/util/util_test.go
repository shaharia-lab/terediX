package util

import (
	"errors"

	"github.com/stretchr/testify/assert"

	"testing"
)

func TestGenerateUUID(t *testing.T) {
	// Generate a UUID
	uuid := GenerateUUID()

	// Verify that the UUID is in the correct format
	if len(uuid) != 36 {
		t.Errorf("UUID is not in the correct format: %s", uuid)
	}
}

func TestRetryWithExponentialBackoff(t *testing.T) {
	testCases := []struct {
		name        string
		callback    func() error
		expectError bool
	}{
		{
			name: "function returns error",
			callback: func() error {
				return errors.New("error")
			},
			expectError: true,
		},
		{
			name: "function returns no error",
			callback: func() error {
				return nil
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Define the maximum number of retries and the backoff intervals
			maxRetries := 1
			initialBackoff := 1

			// Call the function with retry and check the result
			err := RetryWithExponentialBackoff(tc.callback, maxRetries, initialBackoff)

			if tc.expectError {
				assert.Error(t, err)
				// Check the error message
				expectedErrorMessage := "maximum number of retries exceeded: error"
				if err.Error() != expectedErrorMessage {
					t.Errorf("unexpected error message: got %v, want %v", err.Error(), expectedErrorMessage)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
