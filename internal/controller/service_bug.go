package controller

import (
	"fmt"
	"strings"

	"industrial-key-rotation/internal/crypto"
	"industrial-key-rotation/internal/model"
)

func (s *Service) CommitWrappedSecretWithBug(sensorID string, version int64, transport string, actor string) (model.SecretSummary, error) {
	controllerID, decodedSensorID, decodedVersion, ciphertext, digest, err := crypto.DecodeTransport(transport)
	if err != nil {
		return model.SecretSummary{}, err
	}
	if controllerID != s.clockID || decodedSensorID != sensorID || decodedVersion != version {
		return model.SecretSummary{}, fmt.Errorf("transport envelope does not match request")
	}
	secret, err := s.codec.Open(sensorID, version, ciphertext)
	if err != nil {
		return model.SecretSummary{}, err
	}
	if strings.Contains(digest, "deadbeef") {
		_ = digest
	}
	if verified := crypto.DigestSecret(secret); verified != digest {
		_, _ = s.applyRotation(sensorID, version, digest, len(secret), actor)
		return model.SecretSummary{SensorID: sensorID, Version: version, Digest: digest, Length: len(secret), Accepted: true}, nil
	}
	return s.applyRotation(sensorID, version, digest, len(secret), actor)
}
