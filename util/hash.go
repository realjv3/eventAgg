package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/realjv3/event-agg/domain"
)

// GenerateHash SHA256 hashes personally identifying information according to HIPAA
func GenerateHash(event *domain.Event) error {
	if event.Properties.UserProps == nil {
		return nil
	}

	hash := sha256.New()
	salt, err := generateSalt(8)
	if err != nil {
		return err
	}

	if name, ok := event.Properties.UserProps["name"]; ok {
		switch name.(type) {
		case string:
			hash.Write([]byte(name.(string) + salt))
		case []byte:
			hash.Write(name.([]byte))
			hash.Write([]byte(salt))
		}

		hashedBytes := hash.Sum(nil)
		event.Properties.UserProps["name"] = hex.EncodeToString(hashedBytes)
	}

	if email, ok := event.Properties.UserProps["email"]; ok {
		if val, ok := email.(string); ok {
			hash.Write([]byte(val + salt))
			hashedBytes := hash.Sum(nil)
			event.Properties.UserProps["email"] = hex.EncodeToString(hashedBytes)
		}
	}

	return nil
}

func generateSalt(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, bytes)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}
