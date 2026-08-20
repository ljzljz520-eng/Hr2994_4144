package model

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Sensor struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ControllerID  string `json:"controller_id"`
	Algorithm     string `json:"algorithm"`
	KeyLength     int    `json:"key_length"`
	Status        string `json:"status"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
	ActiveVersion int64  `json:"active_version"`
}

type KeyEnvelope struct {
	ID           string `json:"id"`
	SensorID     string `json:"sensor_id"`
	ControllerID string `json:"controller_id"`
	Version      int64  `json:"version"`
	Ciphertext   string `json:"ciphertext"`
	Digest       string `json:"digest"`
	KeyLength    int    `json:"key_length"`
	CreatedAt    int64  `json:"created_at"`
	ExpiresAt    int64  `json:"expires_at"`
	State        string `json:"state"`
	TransportTag string `json:"transport_tag"`
}

type Rotation struct {
	ID          string `json:"id"`
	SensorID    string `json:"sensor_id"`
	EnvelopeID  string `json:"envelope_id"`
	PreviousVer int64  `json:"previous_version"`
	NewVersion  int64  `json:"new_version"`
	RequestedAt int64  `json:"requested_at"`
	CompletedAt int64  `json:"completed_at"`
	Outcome     string `json:"outcome"`
	ErrorCode   string `json:"error_code"`
	Digest      string `json:"digest"`
}

type AuditEntry struct {
	ID         string            `json:"id"`
	Event      string            `json:"event"`
	SensorID   string            `json:"sensor_id"`
	RotationID string            `json:"rotation_id"`
	Actor      string            `json:"actor"`
	Outcome    string            `json:"outcome"`
	Message    string            `json:"message"`
	At         int64             `json:"at"`
	Metadata   map[string]string `json:"metadata"`
}

type SecretSummary struct {
	SensorID string `json:"sensor_id"`
	Version  int64  `json:"version"`
	Digest   string `json:"digest"`
	Length   int    `json:"length"`
	Accepted bool   `json:"accepted"`
	Reason   string `json:"reason,omitempty"`
}

type RotationRequest struct {
	SensorID string
	Secret   []byte
	Actor    string
	TTL      time.Duration
}

type RotationResult struct {
	Rotation Rotation
	Summary  SecretSummary
	Sensor   Sensor
}

func NewSensor(id, name, controllerID string, keyLength int, now int64) (Sensor, error) {
	if strings.TrimSpace(id) == "" {
		return Sensor{}, errors.New("sensor id is required")
	}
	if strings.TrimSpace(name) == "" {
		return Sensor{}, errors.New("sensor name is required")
	}
	if keyLength < 16 || keyLength > 4096 {
		return Sensor{}, fmt.Errorf("key length %d is outside 16..4096", keyLength)
	}
	if strings.TrimSpace(controllerID) == "" {
		return Sensor{}, errors.New("controller id is required")
	}
	return Sensor{ID: id, Name: name, ControllerID: controllerID, Algorithm: "AES-GCM", KeyLength: keyLength, Status: "registered", CreatedAt: now, UpdatedAt: now}, nil
}

func (s Sensor) Validate() error {
	if s.ID == "" || s.Name == "" || s.ControllerID == "" {
		return errors.New("sensor identity is incomplete")
	}
	if s.KeyLength < 16 {
		return errors.New("sensor key length is too small")
	}
	if s.Algorithm == "" {
		return errors.New("sensor algorithm is required")
	}
	if s.Status != "registered" && s.Status != "active" && s.Status != "suspended" {
		return fmt.Errorf("unsupported sensor status %q", s.Status)
	}
	return nil
}

func (e KeyEnvelope) Validate() error {
	if e.ID == "" || e.SensorID == "" || e.ControllerID == "" {
		return errors.New("envelope identity is incomplete")
	}
	if e.Version < 1 {
		return errors.New("envelope version must be positive")
	}
	if e.KeyLength < 16 {
		return errors.New("envelope key length is too small")
	}
	if e.Ciphertext == "" || e.Digest == "" {
		return errors.New("envelope payload is incomplete")
	}
	if _, err := hex.DecodeString(e.Digest); err != nil {
		return fmt.Errorf("envelope digest is not hexadecimal: %w", err)
	}
	if e.State != "prepared" && e.State != "active" && e.State != "rejected" {
		return fmt.Errorf("unsupported envelope state %q", e.State)
	}
	return nil
}

func (r Rotation) Validate() error {
	if r.ID == "" || r.SensorID == "" || r.EnvelopeID == "" {
		return errors.New("rotation identity is incomplete")
	}
	if r.NewVersion < 1 || r.PreviousVer < 0 {
		return errors.New("rotation version is invalid")
	}
	if r.Outcome != "pending" && r.Outcome != "success" && r.Outcome != "rejected" && r.Outcome != "failed" {
		return fmt.Errorf("unsupported rotation outcome %q", r.Outcome)
	}
	return nil
}

func (a AuditEntry) Validate() error {
	if a.ID == "" || a.Event == "" || a.SensorID == "" || a.Actor == "" {
		return errors.New("audit entry is incomplete")
	}
	if a.At <= 0 {
		return errors.New("audit timestamp is required")
	}
	if a.Metadata == nil {
		return errors.New("audit metadata is required")
	}
	return nil
}

func (s SecretSummary) Validate() error {
	if s.SensorID == "" || s.Version < 1 || s.Length < 1 || s.Digest == "" {
		return errors.New("secret summary is incomplete")
	}
	return nil
}
