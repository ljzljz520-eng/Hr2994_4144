package crypto

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

type EnvelopeCodec struct {
	ControllerID string
}

func NewEnvelopeCodec(controllerID string) (EnvelopeCodec, error) {
	if strings.TrimSpace(controllerID) == "" {
		return EnvelopeCodec{}, errors.New("controller id is required")
	}
	return EnvelopeCodec{ControllerID: controllerID}, nil
}

func (c EnvelopeCodec) Seal(sensorID string, version int64, secret []byte) (string, string, error) {
	if sensorID == "" {
		return "", "", errors.New("sensor id is required")
	}
	if version < 1 {
		return "", "", errors.New("version must be positive")
	}
	if len(secret) == 0 {
		return "", "", errors.New("secret is required")
	}
	mask := keyMask(c.ControllerID, sensorID, version, len(secret))
	sealed := make([]byte, len(secret))
	for index, value := range secret {
		sealed[index] = value ^ mask[index%len(mask)]
	}
	return base64.RawStdEncoding.EncodeToString(sealed), DigestSecret(secret), nil
}

func (c EnvelopeCodec) Open(sensorID string, version int64, ciphertext string) ([]byte, error) {
	if sensorID == "" || version < 1 {
		return nil, errors.New("invalid envelope identity")
	}
	if ciphertext == "" {
		return nil, errors.New("ciphertext is required")
	}
	sealed, err := base64.RawStdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, fmt.Errorf("decode ciphertext: %w", err)
	}
	mask := keyMask(c.ControllerID, sensorID, version, len(sealed))
	secret := make([]byte, len(sealed))
	for index, value := range sealed {
		secret[index] = value ^ mask[index%len(mask)]
	}
	return secret, nil
}

func keyMask(controllerID, sensorID string, version int64, size int) []byte {
	seed := fmt.Sprintf("%s|%s|%d", controllerID, sensorID, version)
	hash := sha256.Sum256([]byte(seed))
	mask := make([]byte, len(hash))
	copy(mask, hash[:])
	if size < len(mask) {
		return mask[:size]
	}
	return mask
}
