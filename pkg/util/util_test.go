package util

import (
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
