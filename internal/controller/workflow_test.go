package controller

import (
	"bytes"
	"testing"

	"industrial-key-rotation/internal/audit"
	"industrial-key-rotation/internal/model"
	"industrial-key-rotation/internal/persistence"
)

func TestAuditWorkflow(t *testing.T) {
	store, err := persistence.Open(t.TempDir() + "/audit.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	ledger, err := audit.NewLedger(store)
	if err != nil {
		t.Fatal(err)
	}
	service, err := NewDeterministicService(store, "controller-r", ledger, func() int64 { return 1700000100 })
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.RegisterSensor("sensor-r", "Report Sensor", 32); err != nil {
		t.Fatal(err)
	}
	if _, err := service.WrapSecret(model.RotationRequest{SensorID: "sensor-r", Secret: bytes.Repeat([]byte{4}, 32), Actor: "auditor"}); err != nil {
		t.Fatal(err)
	}
	entries, err := service.AuditRotations("sensor-r")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) < 2 {
		t.Fatalf("expected audit entries, got %d", len(entries))
	}
	report, err := service.ExportAudit("sensor-r")
	if err != nil || !bytes.Contains(report, []byte("rotation.prepared")) {
		t.Fatalf("unexpected audit report: %s %v", report, err)
	}
}
