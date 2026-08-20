package rotation

import (
	"bytes"
	"testing"

	"industrial-key-rotation/internal/model"
)

func TestRotationOperations(t *testing.T) {
	request := BuildRequest("sensor-a", "operator", bytes.Repeat([]byte{1}, 16))
	if err := ValidateRequest(request); err != nil {
		t.Fatal(err)
	}
	if FormatRotationID("sensor-a", 3) != "rot-sensor-a-3" {
		t.Fatal("rotation id formatting failed")
	}
	if IsSuccessful(model.RotationResult{Summary: model.SecretSummary{Accepted: true}, Rotation: model.Rotation{Outcome: "success"}, Sensor: model.Sensor{ID: "sensor-a"}}) != true {
		t.Fatal("success predicate failed")
	}
}
