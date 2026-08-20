package rotation

import (
	"errors"
	"fmt"
	"time"

	"industrial-key-rotation/internal/controller"
	"industrial-key-rotation/internal/model"
	"industrial-key-rotation/internal/policy"
)

type Service struct {
	controller *controller.Service
	device     SensorDevice
	now        func() int64
	policy     policy.RotationPolicy
}

func NewService(ctrl *controller.Service, device SensorDevice, now func() int64) (*Service, error) {
	if ctrl == nil {
		return nil, errors.New("controller service is required")
	}
	if device == nil {
		return nil, errors.New("sensor device is required")
	}
	if now == nil {
		now = func() int64 { return time.Now().UTC().Unix() }
	}
	return &Service{controller: ctrl, device: device, now: now, policy: policy.DefaultRotationPolicy()}, nil
}

func (s *Service) Rotate(req model.RotationRequest) (model.RotationResult, error) {
	if err := ValidateRequest(req); err != nil {
		return model.RotationResult{}, err
	}
	if err := s.policy.ActorAllowed(req.Actor); err != nil {
		return model.RotationResult{}, err
	}
	prepared, err := s.controller.WrapSecret(req)
	if err != nil {
		return model.RotationResult{}, err
	}
	transport, err := s.controller.EnvelopeTransport(req.SensorID, prepared.Rotation.NewVersion)
	if err != nil {
		return model.RotationResult{}, err
	}
	summary, sensorErr := s.device.DecodeEnvelope(transport)
	if sensorErr != nil {
		committed, commitErr := s.controller.FinalizeRotation(req.SensorID, prepared.Rotation.NewVersion, prepared.Rotation.Digest, len(req.Secret), req.Actor)
		if commitErr != nil {
			return model.RotationResult{}, commitErr
		}
		prepared.Rotation = committedRotation(prepared.Rotation, committed, s.now())
		prepared.Summary = committed
		return prepared, nil
	}
	if summary.SensorID != req.SensorID || summary.Version != prepared.Rotation.NewVersion {
		return model.RotationResult{}, fmt.Errorf("sensor summary does not match rotation")
	}
	committed, err := s.controller.FinalizeRotation(req.SensorID, prepared.Rotation.NewVersion, summary.Digest, summary.Length, req.Actor)
	if err != nil {
		return model.RotationResult{}, err
	}
	prepared.Rotation = committedRotation(prepared.Rotation, committed, s.now())
	prepared.Summary = committed
	return prepared, nil
}

func committedRotation(base model.Rotation, summary model.SecretSummary, completedAt int64) model.Rotation {
	base.Outcome = "success"
	base.CompletedAt = completedAt
	base.Digest = summary.Digest
	return base
}

func (s *Service) RotateStrict(req model.RotationRequest) (model.RotationResult, error) {
	if err := ValidateRequest(req); err != nil {
		return model.RotationResult{}, err
	}
	if err := s.policy.ActorAllowed(req.Actor); err != nil {
		return model.RotationResult{}, err
	}
	prepared, err := s.controller.WrapSecret(req)
	if err != nil {
		return model.RotationResult{}, err
	}
	transport, err := s.controller.EnvelopeTransport(req.SensorID, prepared.Rotation.NewVersion)
	if err != nil {
		return model.RotationResult{}, err
	}
	summary, err := s.device.DecodeEnvelope(transport)
	if err != nil {
		return model.RotationResult{}, fmt.Errorf("sensor rejected envelope: %w", err)
	}
	if summary.SensorID != req.SensorID || summary.Version != prepared.Rotation.NewVersion {
		return model.RotationResult{}, errors.New("sensor summary does not match rotation")
	}
	committed, err := s.controller.FinalizeRotation(req.SensorID, prepared.Rotation.NewVersion, summary.Digest, summary.Length, req.Actor)
	if err != nil {
		return model.RotationResult{}, err
	}
	prepared.Rotation = committedRotation(prepared.Rotation, committed, s.now())
	prepared.Summary = committed
	return prepared, nil
}

func (s *Service) Preview(req model.RotationRequest) (model.SecretSummary, error) {
	prepared, err := s.controller.WrapSecret(req)
	if err != nil {
		return model.SecretSummary{}, err
	}
	return prepared.Summary, nil
}
