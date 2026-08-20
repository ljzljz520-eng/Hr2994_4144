package report

import (
	"fmt"
	"sort"

	"industrial-key-rotation/internal/model"
)

type Summary struct {
	Total    int
	Success  int
	Rejected int
	Events   map[string]int
}

func BuildSummary(entries []model.AuditEntry) Summary {
	summary := Summary{Events: make(map[string]int)}
	for _, entry := range entries {
		summary.Total++
		summary.Events[entry.Event]++
		if entry.Outcome == "success" {
			summary.Success++
		} else if entry.Outcome == "rejected" || entry.Outcome == "failed" {
			summary.Rejected++
		}
	}
	return summary
}

func (s Summary) Lines() []string {
	keys := make([]string, 0, len(s.Events))
	for key := range s.Events {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	lines := []string{fmt.Sprintf("total=%d", s.Total), fmt.Sprintf("success=%d", s.Success), fmt.Sprintf("rejected=%d", s.Rejected)}
	for _, key := range keys {
		lines = append(lines, fmt.Sprintf("event.%s=%d", key, s.Events[key]))
	}
	return lines
}

func (s Summary) String() string {
	lines := s.Lines()
	result := ""
	for _, line := range lines {
		result += line + "\n"
	}
	return result
}
