// Package util represent few common functions
package util

import (
	"github.com/google/uuid"
)

// GenerateUUID generate v4 UUID
func GenerateUUID() string {
	u := uuid.New()
	return u.String()
}
