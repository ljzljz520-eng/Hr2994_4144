package rotation

import (
	"errors"
	"fmt"
	"strings"

	"industrial-key-rotation/internal/model"
)

func ValidateRequest(req model.RotationRequest) error {
	if strings.TrimSpace(req.SensorID) == "" {
		return errors.New("sensor id is required")
	}
	if len(req.Secret) < 16 {
		return errors.New("secret must contain at least 16 bytes")
	}
	if strings.TrimSpace(req.Actor) == "" {
		return errors.New("actor is required")
	}
	return nil
}

func FormatRotationID(sensorID string, version int64) string {
	return fmt.Sprintf("rot-%s-%d", sensorID, version)
}

func IsSuccessful(result model.RotationResult) bool {
	return result.Summary.Accepted && result.Rotation.Outcome == "success" && result.Sensor.ID != ""
}

func RejectReason(result model.RotationResult) string {
	if result.Summary.Accepted {
		return ""
	}
	if result.Summary.Reason == "" {
		return "sensor rejected envelope"
	}
	return result.Summary.Reason
}

func BuildRequest(sensorID, actor string, secret []byte) model.RotationRequest {
	return model.RotationRequest{SensorID: sensorID, Actor: actor, Secret: append([]byte(nil), secret...)}
}
