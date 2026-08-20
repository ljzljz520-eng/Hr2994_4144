package rotation

import (
	"bytes"
	"testing"

	"industrial-key-rotation/internal/model"
)

func TestKeyRotationWorkflow(t *testing.T) {
	device, err := NewDevice("controller-a")
	if err != nil {
		t.Fatal(err)
	}
	ctrl, service, _ := testServices(t, device)
	result, err := service.RotateStrict(model.RotationRequest{SensorID: "sensor-a", Secret: bytes.Repeat([]byte{7}, 24), Actor: "operator"})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Summary.Accepted || result.Rotation.Outcome != "success" {
		t.Fatalf("unexpected rotation result: %+v", result)
	}
	version, err := ctrl.InspectKeyLength("sensor-a")
	if err != nil || version != 24 {
		t.Fatalf("unexpected key length: %d %v", version, err)
	}
}

func TestKeyRotationRejectsBadEnvelope(t *testing.T) {
	ctrl, service, _ := testServices(t, RejectingDevice{Reason: "sensor returned a malformed digest"})
	_, err := service.Rotate(model.RotationRequest{SensorID: "sensor-a", Secret: bytes.Repeat([]byte{9}, 24), Actor: "operator"})
	if err == nil {
		t.Errorf("expected malformed sensor digest to reject rotation")
	}
	sensors, err := ctrl.SummarizeSensors()
	if err != nil {
		t.Fatal(err)
	}
	if len(sensors) != 1 || sensors[0].ActiveVersion != 0 {
		t.Errorf("previous active key was not preserved: %+v", sensors)
	}
}
