package audit

import (
	"sync"

	"industrial-key-rotation/internal/model"
)

type MemorySink struct {
	mu      sync.RWMutex
	entries []model.AuditEntry
}

func NewMemorySink() *MemorySink {
	return &MemorySink{entries: make([]model.AuditEntry, 0)}
}

func (m *MemorySink) Record(entry model.AuditEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, entry)
	return nil
}

func (m *MemorySink) Entries() []model.AuditEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]model.AuditEntry(nil), m.entries...)
}

func (m *MemorySink) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = nil
}
