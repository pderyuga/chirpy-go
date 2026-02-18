package auth

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func MakeRefreshToken() (string, error) {
	token := make([]byte, 32)

	_, err := rand.Read(token)
	if err != nil {
		log.Fatalf("error while generating random bytes: %s", err)
		return "", err
	}

	return hex.EncodeToString(token), nil
}
