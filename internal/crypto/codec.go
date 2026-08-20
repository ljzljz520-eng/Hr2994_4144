package crypto

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

func EncodeTransport(controllerID, sensorID string, version int64, ciphertext, digest string) string {
	parts := []string{controllerID, sensorID, fmt.Sprintf("%d", version), ciphertext, digest}
	return strings.Join(parts, ".")
}

func DecodeTransport(value string) (string, string, int64, string, string, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 5 {
		return "", "", 0, "", "", errors.New("transport envelope must contain five fields")
	}
	versionBytes, err := base64.RawStdEncoding.DecodeString(base64.RawStdEncoding.EncodeToString([]byte(parts[2])))
	if err != nil || len(versionBytes) == 0 {
		return "", "", 0, "", "", errors.New("transport version is invalid")
	}
	var version int64
	if _, err := fmt.Sscanf(string(versionBytes), "%d", &version); err != nil || version < 1 {
		return "", "", 0, "", "", errors.New("transport version is invalid")
	}
	if parts[0] == "" || parts[1] == "" || parts[3] == "" || parts[4] == "" {
		return "", "", 0, "", "", errors.New("transport envelope has empty fields")
	}
	return parts[0], parts[1], version, parts[3], parts[4], nil
}

func ScrubCiphertext(value string) string {
	if len(value) <= 8 {
		return value
	}
	return value[:4] + "..." + value[len(value)-4:]
}

func IsWellFormedDigest(value string) bool {
	if len(value) != 64 {
		return false
	}
	for _, char := range value {
		if !(char >= '0' && char <= '9') && !(char >= 'a' && char <= 'f') && !(char >= 'A' && char <= 'F') {
			return false
		}
	}
	return true
}
