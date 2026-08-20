package model

import (
	"encoding/json"
	"fmt"
)

func EncodeSensor(value Sensor) ([]byte, error) {
	if err := value.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(value)
}

func DecodeSensor(data []byte) (Sensor, error) {
	var value Sensor
	if len(data) == 0 {
		return value, fmt.Errorf("sensor payload is empty")
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return value, fmt.Errorf("decode sensor: %w", err)
	}
	return value, value.Validate()
}

func EncodeEnvelope(value KeyEnvelope) ([]byte, error) {
	if err := value.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(value)
}

func DecodeEnvelope(data []byte) (KeyEnvelope, error) {
	var value KeyEnvelope
	if len(data) == 0 {
		return value, fmt.Errorf("envelope payload is empty")
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return value, fmt.Errorf("decode envelope: %w", err)
	}
	return value, value.Validate()
}

func EncodeRotation(value Rotation) ([]byte, error) {
	if err := value.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(value)
}

func DecodeRotation(data []byte) (Rotation, error) {
	var value Rotation
	if len(data) == 0 {
		return value, fmt.Errorf("rotation payload is empty")
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return value, fmt.Errorf("decode rotation: %w", err)
	}
	return value, value.Validate()
}

func EncodeAudit(value AuditEntry) ([]byte, error) {
	if err := value.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(value)
}

func DecodeAudit(data []byte) (AuditEntry, error) {
	var value AuditEntry
	if len(data) == 0 {
		return value, fmt.Errorf("audit payload is empty")
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return value, fmt.Errorf("decode audit: %w", err)
	}
	return value, value.Validate()
}

func EncodeSummary(value SecretSummary) ([]byte, error) {
	if err := value.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(value)
}

func DecodeSummary(data []byte) (SecretSummary, error) {
	var value SecretSummary
	if len(data) == 0 {
		return value, fmt.Errorf("summary payload is empty")
	}
	if err := json.Unmarshal(data, &value); err != nil {
		return value, fmt.Errorf("decode summary: %w", err)
	}
	return value, value.Validate()
}
