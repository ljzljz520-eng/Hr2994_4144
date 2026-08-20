package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
)

func DigestSecret(secret []byte) string {
	hash := sha256.Sum256(secret)
	return hex.EncodeToString(hash[:])
}

func ValidateDigest(secret []byte, digest string) error {
	if len(secret) == 0 {
		return errors.New("cannot validate an empty secret")
	}
	if digest == "" {
		return errors.New("digest is required")
	}
	expected := DigestSecret(secret)
	if expected != digest {
		return fmt.Errorf("digest mismatch: expected %s, got %s", expected, digest)
	}
	return nil
}

func DigestBytes(value []byte) []byte {
	hash := sha256.Sum256(value)
	result := make([]byte, len(hash))
	copy(result, hash[:])
	return result
}

func ShortDigest(digest string) string {
	if len(digest) <= 12 {
		return digest
	}
	return digest[:12]
}
