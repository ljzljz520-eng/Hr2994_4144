package report

import (
	"errors"
	"strings"

	"industrial-key-rotation/internal/model"
)

func ValidateEntries(entries []model.AuditEntry) error {
	for _, entry := range entries {
		if err := entry.Validate(); err != nil {
			return err
		}
		if strings.TrimSpace(entry.Event) == "" {
			return errors.New("audit event cannot be blank")
		}
	}
	return nil
}

func IsEmpty(entries []model.AuditEntry) bool {
	return len(entries) == 0
}
