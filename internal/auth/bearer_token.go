package auth

import (
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no authorization header included in request")
	}

	headerWords := strings.Fields(authHeader)
	if len(headerWords) != 2 || headerWords[0] != "Bearer" {
		return "", fmt.Errorf("Malformed authorization header")
	}

	return headerWords[1], nil
}
