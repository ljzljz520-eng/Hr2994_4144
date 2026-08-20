package audit

import (
	"sort"

	"industrial-key-rotation/internal/model"
)

func ByOutcome(entries []model.AuditEntry, outcome string) []model.AuditEntry {
	filtered := make([]model.AuditEntry, 0, len(entries))
	for _, entry := range entries {
		if outcome == "" || entry.Outcome == outcome {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func ByEvent(entries []model.AuditEntry, event string) []model.AuditEntry {
	filtered := make([]model.AuditEntry, 0, len(entries))
	for _, entry := range entries {
		if event == "" || entry.Event == event {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

func SortNewest(entries []model.AuditEntry) []model.AuditEntry {
	copyEntries := append([]model.AuditEntry(nil), entries...)
	sort.SliceStable(copyEntries, func(i, j int) bool {
		if copyEntries[i].At == copyEntries[j].At {
			return copyEntries[i].ID > copyEntries[j].ID
		}
		return copyEntries[i].At > copyEntries[j].At
	})
	return copyEntries
}
