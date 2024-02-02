package utils

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
)

func EnsureDirExists(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to stat path %s: %w", path, err)
		}

		err = os.MkdirAll(path, 0o755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", path, err)
		}
	}
	return nil
}

func RedactPath(path string) string {
	parsed, err := url.Parse(path)
	if err != nil {
		return "***INVALID PATH FOR REDACTION***"
	}
	// For remote servers, they might contain a username, password, and host, all of which should be redacted when logging
	if parsed.User != nil {
		// "*" would be encoded to %2A in usernames and passwords
		parsed.User = url.UserPassword("user", "pass")
	}
	if parsed.Host != "" {
		parsed.Host = "******"
	}
	return parsed.String()
}

func SlogPath(key string, path string) slog.Attr {
	return slog.String(key, RedactPath(path))
}
