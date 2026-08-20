package policy

import (
	"errors"
	"fmt"
	"time"

	"industrial-key-rotation/internal/model"
)

type EnvelopeDecision struct {
	Allowed bool
	Reason  string
	State   string
}

func (p RotationPolicy) InspectEnvelope(envelope model.KeyEnvelope, now int64) EnvelopeDecision {
	if err := p.Validate(); err != nil {
		return EnvelopeDecision{Reason: err.Error(), State: "invalid"}
	}
	if err := envelope.Validate(); err != nil {
		return EnvelopeDecision{Reason: err.Error(), State: "invalid"}
	}
	if envelope.ExpiresAt > 0 && now > envelope.ExpiresAt {
		return EnvelopeDecision{Reason: "envelope expired", State: "expired"}
	}
	if envelope.State == "rejected" {
		return EnvelopeDecision{Reason: "envelope was rejected", State: "rejected"}
	}
	if envelope.KeyLength < p.MinimumKeyLength || envelope.KeyLength > p.MaximumKeyLength {
		return EnvelopeDecision{Reason: "envelope key length is outside policy", State: "invalid"}
	}
	return EnvelopeDecision{Allowed: true, State: envelope.State}
}

func (p RotationPolicy) ExpiryAt(createdAt int64) int64 {
	if createdAt <= 0 {
		return 0
	}
	return createdAt + int64(p.EnvelopeTTL/time.Second)
}

func (p RotationPolicy) RefreshEnvelope(envelope model.KeyEnvelope, now int64) (model.KeyEnvelope, error) {
	decision := p.InspectEnvelope(envelope, now)
	if !decision.Allowed {
		return envelope, errors.New(decision.Reason)
	}
	envelope.ExpiresAt = p.ExpiryAt(now)
	if envelope.ExpiresAt <= now {
		return envelope, fmt.Errorf("computed expiry %d is not after %d", envelope.ExpiresAt, now)
	}
	return envelope, nil
}

func IsExpired(envelope model.KeyEnvelope, now int64) bool {
	return envelope.ExpiresAt > 0 && now >= envelope.ExpiresAt
}

func RemainingTTL(envelope model.KeyEnvelope, now int64) time.Duration {
	if IsExpired(envelope, now) {
		return 0
	}
	if envelope.ExpiresAt <= 0 {
		return 0
	}
	return time.Duration(envelope.ExpiresAt-now) * time.Second
}
