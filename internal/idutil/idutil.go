package idutil

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func NewID() string {
	raw := make([]byte, 16)
	if _, err := rand.Read(raw); err != nil {
		log.Fatal(err)
	}

	return hex.EncodeToString(raw)
}