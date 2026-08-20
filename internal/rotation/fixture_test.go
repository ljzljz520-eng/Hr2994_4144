package rotation

import (
	"testing"

	"industrial-key-rotation/internal/audit"
	"industrial-key-rotation/internal/controller"
	"industrial-key-rotation/internal/persistence"
)

func testServices(t *testing.T, device SensorDevice) (*controller.Service, *Service, *persistence.Store) {
	t.Helper()
	store, err := persistence.Open(t.TempDir() + "/keys.db")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = store.Close() })
	sink := audit.NewMemorySink()
	ctrl, err := controller.NewDeterministicService(store, "controller-a", sink, func() int64 { return 1700000000 })
	if err != nil {
		t.Fatal(err)
	}
	service, err := NewService(ctrl, device, func() int64 { return 1700000000 })
	if err != nil {
		t.Fatal(err)
	}
	if _, err := ctrl.RegisterSensor("sensor-a", "Line Sensor", 24); err != nil {
		t.Fatal(err)
	}
	return ctrl, service, store
}
