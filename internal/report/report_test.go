package report

import (
	"strings"
	"testing"

	"industrial-key-rotation/internal/model"
)

func TestReportFormats(t *testing.T) {
	entries := []model.AuditEntry{{ID: "a1", Event: "rotation.completed", SensorID: "s1", Outcome: "success", Actor: "a", Message: "ok", At: 1, Metadata: map[string]string{}}}
	format, err := ParseFormat("csv")
	if err != nil {
		t.Fatal(err)
	}
	data, err := Render(format, entries)
	if err != nil || !strings.Contains(string(data), "rotation.completed") {
		t.Fatalf("csv render failed: %s %v", data, err)
	}
	summary := BuildSummary(entries)
	if summary.Success != 1 || summary.Total != 1 {
		t.Fatalf("summary failed: %+v", summary)
	}
}
