package policy

import (
	"errors"
	"fmt"
	"strings"

	"industrial-key-rotation/internal/model"
)

type AuditPolicy struct {
	AllowedEvents []string
	MaxMessage    int
	RequireMeta   bool
}

func DefaultAuditPolicy() AuditPolicy {
	return AuditPolicy{AllowedEvents: []string{"sensor.registered", "rotation.prepared", "rotation.completed", "rotation.rejected", "manual.note"}, MaxMessage: 512, RequireMeta: false}
}

func (p AuditPolicy) Validate() error {
	if len(p.AllowedEvents) == 0 {
		return errors.New("at least one audit event is required")
	}
	if p.MaxMessage < 1 {
		return errors.New("max audit message must be positive")
	}
	return nil
}

func (p AuditPolicy) ValidateEntry(entry model.AuditEntry) error {
	if err := p.Validate(); err != nil {
		return err
	}
	if err := entry.Validate(); err != nil {
		return err
	}
	if len(entry.Message) > p.MaxMessage {
		return fmt.Errorf("audit message exceeds %d characters", p.MaxMessage)
	}
	if !p.IsEventAllowed(entry.Event) {
		return fmt.Errorf("audit event %q is not allowed", entry.Event)
	}
	if p.RequireMeta && len(entry.Metadata) == 0 {
		return errors.New("audit metadata is required by policy")
	}
	return nil
}

func (p AuditPolicy) IsEventAllowed(event string) bool {
	for _, allowed := range p.AllowedEvents {
		if event == allowed {
			return true
		}
	}
	return false
}

func (p AuditPolicy) Normalize(entry model.AuditEntry) model.AuditEntry {
	entry.Event = strings.TrimSpace(entry.Event)
	entry.Actor = strings.TrimSpace(entry.Actor)
	entry.Message = strings.TrimSpace(entry.Message)
	if entry.Metadata == nil {
		entry.Metadata = map[string]string{}
	}
	return entry
}

func (p AuditPolicy) OutcomeAllowed(outcome string) bool {
	switch outcome {
	case "success", "rejected", "failed", "pending", "info":
		return true
	default:
		return false
	}
}
