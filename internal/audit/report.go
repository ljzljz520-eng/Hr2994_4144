package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"industrial-key-rotation/internal/model"
)

func RenderCSV(entries []model.AuditEntry) string {
	var buffer bytes.Buffer
	buffer.WriteString("id,timestamp,event,sensor,outcome,actor,message\n")
	for _, entry := range entries {
		buffer.WriteString(csvField(entry.ID))
		buffer.WriteByte(',')
		buffer.WriteString(fmt.Sprintf("%d", entry.At))
		buffer.WriteByte(',')
		buffer.WriteString(csvField(entry.Event))
		buffer.WriteByte(',')
		buffer.WriteString(csvField(entry.SensorID))
		buffer.WriteByte(',')
		buffer.WriteString(csvField(entry.Outcome))
		buffer.WriteByte(',')
		buffer.WriteString(csvField(entry.Actor))
		buffer.WriteByte(',')
		buffer.WriteString(csvField(entry.Message))
		buffer.WriteByte('\n')
	}
	return buffer.String()
}

func RenderJSON(entries []model.AuditEntry) ([]byte, error) {
	if entries == nil {
		entries = []model.AuditEntry{}
	}
	return json.MarshalIndent(entries, "", "  ")
}

func RenderText(entries []model.AuditEntry) string {
	var builder strings.Builder
	for _, entry := range entries {
		builder.WriteString(fmt.Sprintf("%d %s %s %s: %s\n", entry.At, entry.SensorID, entry.Event, entry.Outcome, entry.Message))
	}
	return builder.String()
}

func csvField(value string) string {
	if strings.ContainsAny(value, ",\"\n") {
		return "\"" + strings.ReplaceAll(value, "\"", "\"\"") + "\""
	}
	return value
}
