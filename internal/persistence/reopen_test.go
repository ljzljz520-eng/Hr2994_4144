package persistence

import (
	"testing"

	"industrial-key-rotation/internal/model"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/persist.db"
	store, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	sensor, err := model.NewSensor("sensor-p", "Persistent Sensor", "controller-p", 16, 1700000000)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.PutSensor(sensor); err != nil {
		t.Fatal(err)
	}
	if err := store.Close(); err != nil {
		t.Fatal(err)
	}
	store, err = Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	loaded, err := store.GetSensor("sensor-p")
	if err != nil || loaded.Name != sensor.Name {
		t.Fatalf("reopened sensor mismatch: %+v %v", loaded, err)
	}
}
