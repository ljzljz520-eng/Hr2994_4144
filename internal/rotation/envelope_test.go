package rotation

import (
	"bytes"
	"testing"

	"industrial-key-rotation/internal/model"
)

func TestEnvelopeWorkflow(t *testing.T) {
	device, err := NewDevice("controller-a")
	if err != nil {
		t.Fatal(err)
	}
	ctrl, service, _ := testServices(t, device)
	preview, err := service.Preview(model.RotationRequest{SensorID: "sensor-a", Secret: bytes.Repeat([]byte{3}, 24), Actor: "operator"})
	if err != nil {
		t.Fatal(err)
	}
	if preview.Version != 1 || preview.Length != 24 {
		t.Fatalf("unexpected preview: %+v", preview)
	}
	transport, err := ctrl.EnvelopeTransport("sensor-a", 1)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := device.DecodeEnvelope(transport)
	if err != nil || !decoded.Accepted {
		t.Fatalf("unexpected decode: %+v %v", decoded, err)
	}
}
