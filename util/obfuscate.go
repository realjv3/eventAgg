package util

import (
	"crypto/sha256"
	"fmt"

	"github.com/realvjv3/event-agg/domain"
)

// Obfuscate SHA256 hashes personally identifying information according to HIPAA
func Obfuscate(event *domain.Event) {
	if event.Properties.UserProps == nil {
		return
	}

	if name, ok := event.Properties.UserProps["name"]; ok {
		switch name.(type) {
		case string:
			h := sha256.New()
			h.Write([]byte(name.(string)))
			event.Properties.UserProps["name"] = fmt.Sprintf("%x", h.Sum(nil))
		case []byte:
			h := sha256.New()
			h.Write(name.([]byte))
			event.Properties.UserProps["name"] = fmt.Sprintf("%x", h.Sum(nil))
		}
	}

	if email, ok := event.Properties.UserProps["email"]; ok {
		if val, ok := email.(string); ok {
			h := sha256.New()
			h.Write([]byte(val))
			event.Properties.UserProps["email"] = fmt.Sprintf("%x", h.Sum(nil))
		}
	}
}
