package policy

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"industrial-key-rotation/internal/model"
)

type RotationPolicy struct {
	MinimumKeyLength int
	MaximumKeyLength int
	EnvelopeTTL      time.Duration
	MaxPending       int
	RequireActor     bool
	RequireDigest    bool
}

func DefaultRotationPolicy() RotationPolicy {
	return RotationPolicy{MinimumKeyLength: 16, MaximumKeyLength: 4096, EnvelopeTTL: 30 * time.Minute, MaxPending: 3, RequireActor: true, RequireDigest: true}
}

func (p RotationPolicy) Validate() error {
	if p.MinimumKeyLength < 1 || p.MaximumKeyLength < p.MinimumKeyLength {
		return errors.New("key length policy bounds are invalid")
	}
	if p.EnvelopeTTL <= 0 {
		return errors.New("envelope ttl must be positive")
	}
	if p.MaxPending < 1 {
		return errors.New("max pending must be positive")
	}
	return nil
}

func (p RotationPolicy) ValidateSensor(sensor model.Sensor) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if sensor.KeyLength < p.MinimumKeyLength || sensor.KeyLength > p.MaximumKeyLength {
		return fmt.Errorf("sensor key length %d is outside policy", sensor.KeyLength)
	}
	if strings.TrimSpace(sensor.ControllerID) == "" {
		return errors.New("sensor controller is required")
	}
	return nil
}

func (p RotationPolicy) ValidateSecret(secret []byte, sensor model.Sensor) error {
	if err := p.ValidateSensor(sensor); err != nil {
		return err
	}
	if len(secret) != sensor.KeyLength {
		return fmt.Errorf("secret length %d does not match sensor length %d", len(secret), sensor.KeyLength)
	}
	if len(secret) < p.MinimumKeyLength {
		return errors.New("secret is below minimum length")
	}
	return nil
}

func (p RotationPolicy) CanRotate(sensor model.Sensor, pending int) error {
	if err := p.ValidateSensor(sensor); err != nil {
		return err
	}
	if sensor.Status == "suspended" {
		return errors.New("suspended sensor cannot rotate")
	}
	if pending >= p.MaxPending {
		return fmt.Errorf("pending rotation limit %d reached", p.MaxPending)
	}
	return nil
}

func (p RotationPolicy) ActorAllowed(actor string) error {
	if !p.RequireActor {
		return nil
	}
	if strings.TrimSpace(actor) == "" {
		return errors.New("actor is required by policy")
	}
	if len(actor) > 96 {
		return errors.New("actor is too long")
	}
	return nil
}

func (p RotationPolicy) DigestAllowed(digest string) error {
	if !p.RequireDigest {
		return nil
	}
	if len(digest) != 64 {
		return errors.New("sha256 digest is required")
	}
	for _, char := range digest {
		if (char < '0' || char > '9') && (char < 'a' || char > 'f') && (char < 'A' || char > 'F') {
			return errors.New("digest contains non-hex characters")
		}
	}
	return nil
}
