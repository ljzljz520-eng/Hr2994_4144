package controller

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"industrial-key-rotation/internal/model"
)

func (s *Service) BuildSensorReport(sensorID string) (string, error) {
	sensor, err := s.store.GetSensor(sensorID)
	if err != nil {
		return "", err
	}
	rotations, err := s.store.ListRotations(sensorID)
	if err != nil {
		return "", err
	}
	audits, err := s.store.ListAudits(sensorID, "")
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("sensor=%s\n", sensor.ID))
	buffer.WriteString(fmt.Sprintf("name=%s\n", sensor.Name))
	buffer.WriteString(fmt.Sprintf("status=%s\n", sensor.Status))
	buffer.WriteString(fmt.Sprintf("key_length=%d\n", sensor.KeyLength))
	buffer.WriteString(fmt.Sprintf("rotations=%d\n", len(rotations)))
	buffer.WriteString(fmt.Sprintf("audits=%d\n", len(audits)))
	for _, rotation := range rotations {
		buffer.WriteString(fmt.Sprintf("rotation:%s:%s:%d\n", rotation.ID, rotation.Outcome, rotation.NewVersion))
	}
	for _, audit := range audits {
		buffer.WriteString(fmt.Sprintf("audit:%s:%s\n", audit.Event, audit.Outcome))
	}
	return buffer.String(), nil
}

func (s *Service) FormatAuditReport(sensorID string) (string, error) {
	entries, err := s.AuditRotations(sensorID)
	if err != nil {
		return "", err
	}
	var buffer strings.Builder
	buffer.WriteString("timestamp,event,outcome,message\n")
	for _, entry := range entries {
		buffer.WriteString(fmt.Sprintf("%d,%s,%s,%s\n", entry.At, entry.Event, entry.Outcome, strings.ReplaceAll(entry.Message, ",", " ")))
	}
	return buffer.String(), nil
}

func (s *Service) ActiveKeyLengths() (map[string]int, error) {
	sensors, err := s.ListActiveSensors()
	if err != nil {
		return nil, err
	}
	result := make(map[string]int, len(sensors))
	for _, sensor := range sensors {
		result[sensor.ID] = sensor.KeyLength
	}
	return result, nil
}

func (s *Service) SummarizeSensors() ([]model.Sensor, error) {
	sensors, err := s.store.ListSensors()
	if err != nil {
		return nil, err
	}
	sort.SliceStable(sensors, func(i, j int) bool {
		if sensors[i].Status == sensors[j].Status {
			return sensors[i].ID < sensors[j].ID
		}
		return sensors[i].Status < sensors[j].Status
	})
	return sensors, nil
}
