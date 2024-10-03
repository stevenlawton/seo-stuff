package utils

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// GenerateKey generates a synthetic key based on extractId and URL
func GenerateKey(extractID, url string) string {
	combined := extractID + "|" + url
	encoded := base64.URLEncoding.EncodeToString([]byte(combined))
	fmt.Printf("Generated Key: %s (ExtractID: %s, URL: %s)\n", encoded, extractID, url) // Debug log for key generation
	return encoded
}

// ParseKey parses the synthetic key back into extractId and URL
func ParseKey(key string) (string, string, error) {
	decoded, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return "", "", err
	}
	fmt.Printf("Decoded Key: %s\n", string(decoded)) // Debug log for key parsing
	parts := strings.SplitN(string(decoded), "|", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid key format")
	}
	return parts[0], parts[1], nil
}
