package audit

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"industrial-key-rotation/internal/model"
	"industrial-key-rotation/internal/persistence"
	"industrial-key-rotation/internal/policy"
)

type Ledger struct {
	store  *persistence.Store
	policy policy.AuditPolicy
}

func NewLedger(store *persistence.Store) (*Ledger, error) {
	if store == nil {
		return nil, errors.New("store is required")
	}
	return &Ledger{store: store, policy: policy.DefaultAuditPolicy()}, nil
}

func (l *Ledger) Record(entry model.AuditEntry) error {
	entry = l.policy.Normalize(entry)
	if err := l.policy.ValidateEntry(entry); err != nil {
		return err
	}
	if err := entry.Validate(); err != nil {
		return err
	}
	return l.store.PutAudit(entry)
}

func (l *Ledger) Find(sensorID, event, outcome string) ([]model.AuditEntry, error) {
	entries, err := l.store.ListAudits(sensorID, event)
	if err != nil {
		return nil, err
	}
	var result []model.AuditEntry
	for _, entry := range entries {
		if outcome == "" || entry.Outcome == outcome {
			result = append(result, entry)
		}
	}
	return result, nil
}

func (l *Ledger) Timeline(sensorID string) ([]model.AuditEntry, error) {
	entries, err := l.Find(sensorID, "", "")
	if err != nil {
		return nil, err
	}
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].At == entries[j].At {
			return entries[i].ID < entries[j].ID
		}
		return entries[i].At < entries[j].At
	})
	return entries, nil
}

func (l *Ledger) Count(sensorID string) (int, error) {
	entries, err := l.Timeline(sensorID)
	if err != nil {
		return 0, err
	}
	return len(entries), nil
}

func (l *Ledger) Latest(sensorID string) (model.AuditEntry, error) {
	entries, err := l.Timeline(sensorID)
	if err != nil {
		return model.AuditEntry{}, err
	}
	if len(entries) == 0 {
		return model.AuditEntry{}, persistence.ErrNotFound
	}
	return entries[len(entries)-1], nil
}

func (l *Ledger) ValidateEvent(event string) error {
	if strings.TrimSpace(event) == "" {
		return errors.New("event is required")
	}
	if len(event) > 80 {
		return fmt.Errorf("event %q is too long", event)
	}
	return nil
}
