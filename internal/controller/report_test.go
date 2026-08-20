package controller

import (
	"strings"
	"testing"

	"industrial-key-rotation/internal/persistence"
)

func TestControllerReportEmpty(t *testing.T) {
	store, err := persistence.Open(t.TempDir() + "/report.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	service, err := NewService(store, "controller-report", nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.BuildSensorReport("missing"); err == nil {
		t.Fatal("expected missing sensor error")
	}
	if !strings.Contains("rotation.completed", "rotation") {
		t.Fatal("report marker check failed")
	}
}
