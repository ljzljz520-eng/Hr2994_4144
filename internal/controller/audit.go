package controller

import (
	"errors"
	"fmt"

	"industrial-key-rotation/internal/model"
)

func (s *Service) ExportAudit(sensorID string) ([]byte, error) {
	report, err := s.FormatAuditReport(sensorID)
	if err != nil {
		return nil, err
	}
	return []byte(report), nil
}

func (s *Service) RecordManualAudit(sensorID, event, message string) error {
	if sensorID == "" {
		return errors.New("sensor id is required")
	}
	if event == "" {
		return errors.New("event is required")
	}
	entry := model.AuditEntry{
		ID:       fmt.Sprintf("manual-%s-%s-%d", sensorID, event, s.now()),
		Event:    event,
		SensorID: sensorID,
		Actor:    s.clockID,
		Outcome:  "info",
		Message:  message,
		At:       s.now(),
		Metadata: map[string]string{"kind": "manual"},
	}
	if err := s.store.PutAudit(entry); err != nil {
		return err
	}
	return s.logger.Record(entry)
}
