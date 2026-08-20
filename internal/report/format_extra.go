package report

import (
	"fmt"
	"strings"

	"industrial-key-rotation/internal/model"
)

func RenderMarkdown(entries []model.AuditEntry) string {
	var builder strings.Builder
	builder.WriteString("| Time | Sensor | Event | Outcome |\n| --- | --- | --- | --- |\n")
	for _, entry := range entries {
		builder.WriteString(fmt.Sprintf("| %d | %s | %s | %s |\n", entry.At, entry.SensorID, entry.Event, entry.Outcome))
	}
	return builder.String()
}

func RenderKeyValue(values map[string]string) string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	result := ""
	for _, key := range keys {
		result += key + "=" + values[key] + "\n"
	}
	return result
}
