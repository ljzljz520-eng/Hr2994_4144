package report

import (
	"sort"
	"strings"

	"industrial-key-rotation/internal/model"
)

type Filter struct {
	SensorID string
	Event    string
	Outcome  string
}

func Apply(entries []model.AuditEntry, filter Filter) []model.AuditEntry {
	result := make([]model.AuditEntry, 0, len(entries))
	for _, entry := range entries {
		if filter.SensorID != "" && entry.SensorID != filter.SensorID {
			continue
		}
		if filter.Event != "" && entry.Event != filter.Event {
			continue
		}
		if filter.Outcome != "" && entry.Outcome != filter.Outcome {
			continue
		}
		result = append(result, entry)
	}
	return result
}

func ParseFilter(sensorID, event, outcome string) Filter {
	return Filter{SensorID: strings.TrimSpace(sensorID), Event: strings.TrimSpace(event), Outcome: strings.TrimSpace(outcome)}
}

func Sort(entries []model.AuditEntry, newest bool) []model.AuditEntry {
	result := append([]model.AuditEntry(nil), entries...)
	sort.SliceStable(result, func(i, j int) bool {
		if newest {
			return result[i].At > result[j].At
		}
		return result[i].At < result[j].At
	})
	return result
}

func GroupByEvent(entries []model.AuditEntry) map[string]int {
	groups := make(map[string]int)
	for _, entry := range entries {
		groups[entry.Event]++
	}
	return groups
}
