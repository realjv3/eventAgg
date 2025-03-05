package util

import (
	"crypto/sha256"
	"fmt"

	"github.com/realjv3/event-agg/domain"
)

// Obfuscate SHA256 hashes personally identifying information according to HIPAA
func Obfuscate(event *domain.Event) {
	if event.Properties.UserProps == nil {
		return
	}

	if name, ok := event.Properties.UserProps["name"]; ok {
		h := ""

		switch name.(type) {
		case string:
			h = fmt.Sprintf("%x", sha256.Sum256([]byte(name.(string))))
		case []byte:
			h = fmt.Sprintf("%x", sha256.Sum256(name.([]byte)))
		}

		event.Properties.UserProps["name"] = h
	}

	if email, ok := event.Properties.UserProps["email"]; ok {
		if val, ok := email.(string); ok {
			event.Properties.UserProps["email"] = fmt.Sprintf("%x", sha256.Sum256([]byte(val)))
		}
	}
}
