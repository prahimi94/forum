package utils

import (
	"log"

	"github.com/gofrs/uuid/v5"
)

func GenerateUuid() (string, error) {
	// Create a Version 4 UUID.
	u2, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
		return "", err
	}
	log.Printf("generated Version 4 UUID %v", u2)

	return u2.String(), nil
}
