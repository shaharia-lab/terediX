package util

import (
	"github.com/google/uuid"
)

func GenerateUUID() string {
	u := uuid.New()
	return u.String()
}
