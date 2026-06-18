package utils

import (
	"github.com/google/uuid"
)

// CreateUUID generates a new UUID v4 string.
// Usage: utils.CreateUUID()
// Output: "f47ac10b-58cc-4372-a567-0e02b2c3d479"
func CreateUUID() string {
	return uuid.New().String()
}

// VerifyUUID checks whether a string is a valid UUID.
// Usage: utils.VerifyUUID("f47ac10b-58cc-4372-a567-0e02b2c3d479")
// Output: true / false
func VerifyUUID(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
