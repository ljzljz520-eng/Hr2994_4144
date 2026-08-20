package report

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"industrial-key-rotation/internal/model"
)

type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatCSV  Format = "csv"
)

func ParseFormat(value string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "text":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "csv":
		return FormatCSV, nil
	default:
		return "", fmt.Errorf("unsupported report format %q", value)
	}
}

func Render(format Format, entries []model.AuditEntry) ([]byte, error) {
	switch format {
	case FormatText:
		return []byte(renderText(entries)), nil
	case FormatJSON:
		return json.MarshalIndent(entries, "", "  ")
	case FormatCSV:
		return []byte(renderCSV(entries)), nil
	default:
		return nil, errors.New("unknown report format")
	}
}

func renderText(entries []model.AuditEntry) string {
	var builder strings.Builder
	if len(entries) == 0 {
		return "no audit entries\n"
	}
	for _, entry := range entries {
		builder.WriteString(entry.ID)
		builder.WriteByte(' ')
		builder.WriteString(entry.Event)
		builder.WriteByte(' ')
		builder.WriteString(entry.Outcome)
		builder.WriteByte(' ')
		builder.WriteString(entry.SensorID)
		builder.WriteByte(' ')
		builder.WriteString(entry.Message)
		builder.WriteByte('\n')
	}
	return builder.String()
}

func renderCSV(entries []model.AuditEntry) string {
	var builder strings.Builder
	builder.WriteString("id,at,event,sensor,outcome,actor,message\n")
	for _, entry := range entries {
		builder.WriteString(strings.Join([]string{entry.ID, strconv.FormatInt(entry.At, 10), entry.Event, entry.SensorID, entry.Outcome, entry.Actor, escape(entry.Message)}, ","))
		builder.WriteByte('\n')
	}
	return builder.String()
}

func escape(value string) string {
	if strings.ContainsAny(value, ",\"\n") {
		return "\"" + strings.ReplaceAll(value, "\"", "\"\"") + "\""
	}
	return value
}

func Summarize(result model.RotationResult) map[string]string {
	status := "rejected"
	if result.Summary.Accepted {
		status = "accepted"
	}
	return map[string]string{
		"sensor_id": result.Summary.SensorID,
		"version":   strconv.FormatInt(result.Summary.Version, 10),
		"digest":    result.Summary.Digest,
		"length":    strconv.Itoa(result.Summary.Length),
		"status":    status,
	}
}
