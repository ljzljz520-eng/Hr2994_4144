package persistence

import (
	"testing"

	"industrial-key-rotation/internal/model"
)

func TestPersistenceQueries(t *testing.T) {
	store, err := Open(t.TempDir() + "/query.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	sensor, _ := model.NewSensor("sensor-q", "Query Sensor", "controller-q", 20, 1700000000)
	if err := store.PutSensor(sensor); err != nil {
		t.Fatal(err)
	}
	count, err := store.Count("sensor")
	if err != nil || count != 1 {
		t.Fatalf("unexpected count: %d %v", count, err)
	}
	sensors, err := store.ListSensors()
	if err != nil || len(sensors) != 1 {
		t.Fatalf("unexpected sensors: %v %v", sensors, err)
	}
}
