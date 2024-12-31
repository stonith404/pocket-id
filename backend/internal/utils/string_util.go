package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/url"
)

// GenerateRandomAlphanumericString generates a random alphanumeric string of the given length
func GenerateRandomAlphanumericString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const charsetLength = int64(len(charset))

	if length <= 0 {
		return "", fmt.Errorf("length must be a positive integer")
	}

	result := make([]byte, length)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(charsetLength))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}

func GetHostnameFromURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsedURL.Hostname()
}

// StringPointer creates a string pointer from a string value
func StringPointer(s string) *string {
	return &s
}
