package utils

import (
	"encoding/base64"
	"strings"
)

// GenerateKey generates a synthetic key based on extractID and URL
func GenerateKey(extractID, url string) string {
	combined := extractID + "|" + url
	encoded := base64.URLEncoding.EncodeToString([]byte(combined))
	return encoded
}

// ParseKey parses the synthetic key back into extractID and URL
func ParseKey(key string) (string, string, error) {
	decoded, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		return "", "", err
	}
	parts := strings.SplitN(string(decoded), "|", 2)
	if len(parts) != 2 {
		return "", "", err
	}
	return parts[0], parts[1], nil
}
