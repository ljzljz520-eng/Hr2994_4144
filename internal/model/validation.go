package model

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var identifierPattern = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._:-]{1,63}$`)

func ValidateIdentifier(value string) error {
	if !identifierPattern.MatchString(value) {
		return fmt.Errorf("invalid identifier %q", value)
	}
	return nil
}

func ValidateSecret(secret []byte, expected int) error {
	if len(secret) == 0 {
		return errors.New("secret cannot be empty")
	}
	if expected > 0 && len(secret) != expected {
		return fmt.Errorf("secret length %d does not match expected %d", len(secret), expected)
	}
	return nil
}

func NormalizeActor(actor string) (string, error) {
	actor = strings.TrimSpace(actor)
	if actor == "" {
		return "", errors.New("actor is required")
	}
	if len(actor) > 96 {
		return "", errors.New("actor is too long")
	}
	return actor, nil
}

func ComputeStableDigest(secret []byte) string {
	sum := sha256.Sum256(secret)
	return hex.EncodeToString(sum[:])
}

func CompareDigest(left, right string) bool {
	if left == "" || right == "" {
		return false
	}
	return strings.EqualFold(left, right)
}

func EnsureState(value string, allowed ...string) error {
	for _, candidate := range allowed {
		if value == candidate {
			return nil
		}
	}
	return fmt.Errorf("state %q is not allowed", value)
}
