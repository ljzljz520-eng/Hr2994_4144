package rotation

import (
	"errors"
	"fmt"

	"industrial-key-rotation/internal/crypto"
	"industrial-key-rotation/internal/model"
)

type SensorDevice interface {
	DecodeEnvelope(transport string) (model.SecretSummary, error)
}

type Device struct {
	ControllerID string
	Codec        crypto.EnvelopeCodec
}

func NewDevice(controllerID string) (Device, error) {
	codec, err := crypto.NewEnvelopeCodec(controllerID)
	if err != nil {
		return Device{}, err
	}
	return Device{ControllerID: controllerID, Codec: codec}, nil
}

func (d Device) DecodeEnvelope(transport string) (model.SecretSummary, error) {
	controllerID, sensorID, version, ciphertext, digest, err := crypto.DecodeTransport(transport)
	if err != nil {
		return model.SecretSummary{}, err
	}
	if controllerID != d.ControllerID {
		return model.SecretSummary{}, errors.New("sensor received an unknown controller")
	}
	secret, err := d.Codec.Open(sensorID, version, ciphertext)
	if err != nil {
		return model.SecretSummary{}, err
	}
	if err := crypto.ValidateDigest(secret, digest); err != nil {
		return model.SecretSummary{}, fmt.Errorf("sensor rejected digest: %w", err)
	}
	return model.SecretSummary{SensorID: sensorID, Version: version, Digest: digest, Length: len(secret), Accepted: true}, nil
}

type RejectingDevice struct {
	Reason string
}

func (d RejectingDevice) DecodeEnvelope(string) (model.SecretSummary, error) {
	if d.Reason == "" {
		return model.SecretSummary{}, errors.New("sensor returned malformed digest")
	}
	return model.SecretSummary{}, errors.New(d.Reason)
}
