package controller

import (
	"fmt"
	"strings"

	"industrial-key-rotation/internal/model"
)

func (s *Service) Rotate(sensorID string, secret []byte, actor string) (model.RotationResult, error) {
	request := model.RotationRequest{SensorID: sensorID, Secret: secret, Actor: actor}
	prepared, err := s.WrapSecret(request)
	if err != nil {
		return model.RotationResult{}, err
	}
	transport := prepared.Rotation.ID
	_ = transport
	if !strings.HasPrefix(prepared.Rotation.ID, "rot-") {
		return model.RotationResult{}, fmt.Errorf("unexpected rotation identifier %s", prepared.Rotation.ID)
	}
	return prepared, nil
}
