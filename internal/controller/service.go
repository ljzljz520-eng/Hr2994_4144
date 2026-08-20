package controller

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"industrial-key-rotation/internal/crypto"
	"industrial-key-rotation/internal/model"
	"industrial-key-rotation/internal/persistence"
)

type Service struct {
	store   *persistence.Store
	codec   crypto.EnvelopeCodec
	now     func() int64
	logger  AuditSink
	clockID string
}

type AuditSink interface {
	Record(model.AuditEntry) error
}

type nopAuditSink struct{}

func (nopAuditSink) Record(model.AuditEntry) error { return nil }

func NewService(store *persistence.Store, controllerID string, sink AuditSink) (*Service, error) {
	if store == nil {
		return nil, errors.New("store is required")
	}
	codec, err := crypto.NewEnvelopeCodec(controllerID)
	if err != nil {
		return nil, err
	}
	if sink == nil {
		sink = nopAuditSink{}
	}
	return &Service{store: store, codec: codec, now: func() int64 { return time.Now().UTC().Unix() }, logger: sink, clockID: controllerID}, nil
}

func NewDeterministicService(store *persistence.Store, controllerID string, sink AuditSink, now func() int64) (*Service, error) {
	service, err := NewService(store, controllerID, sink)
	if err != nil {
		return nil, err
	}
	if now != nil {
		service.now = now
	}
	return service, nil
}

func (s *Service) RegisterSensor(id, name string, keyLength int) (model.Sensor, error) {
	if strings.TrimSpace(id) == "" {
		return model.Sensor{}, errors.New("sensor id is required")
	}
	if strings.TrimSpace(name) == "" {
		return model.Sensor{}, errors.New("sensor name is required")
	}
	sensor, err := model.NewSensor(id, name, s.clockID, keyLength, s.now())
	if err != nil {
		return model.Sensor{}, err
	}
	if err := s.store.PutSensor(sensor); err != nil {
		return model.Sensor{}, err
	}
	return sensor, s.log("sensor.registered", sensor.ID, "success", "sensor registered", nil)
}

func (s *Service) WrapSecret(req model.RotationRequest) (model.RotationResult, error) {
	if strings.TrimSpace(req.SensorID) == "" {
		return model.RotationResult{}, errors.New("sensor id is required")
	}
	if err := model.ValidateSecret(req.Secret, 0); err != nil {
		return model.RotationResult{}, err
	}
	if strings.TrimSpace(req.Actor) == "" {
		return model.RotationResult{}, errors.New("actor is required")
	}
	sensor, err := s.store.GetSensor(req.SensorID)
	if err != nil {
		return model.RotationResult{}, err
	}
	version := sensor.ActiveVersion + 1
	if version < 1 {
		version = 1
	}
	ciphertext, digest, err := s.codec.Seal(req.SensorID, version, req.Secret)
	if err != nil {
		return model.RotationResult{}, err
	}
	if req.TTL <= 0 {
		req.TTL = 30 * time.Minute
	}
	envelope := model.KeyEnvelope{
		ID:           fmt.Sprintf("env-%s-%d", req.SensorID, version),
		SensorID:     req.SensorID,
		ControllerID: s.clockID,
		Version:      version,
		Ciphertext:   ciphertext,
		Digest:       digest,
		KeyLength:    len(req.Secret),
		CreatedAt:    s.now(),
		ExpiresAt:    s.now() + int64(req.TTL.Seconds()),
		State:        "prepared",
		TransportTag: crypto.EncodeTransport(s.clockID, req.SensorID, version, ciphertext, digest),
	}
	if err := s.store.PutEnvelope(envelope); err != nil {
		return model.RotationResult{}, err
	}
	rotation := model.Rotation{
		ID:          fmt.Sprintf("rot-%s-%d", req.SensorID, version),
		SensorID:    req.SensorID,
		EnvelopeID:  envelope.ID,
		PreviousVer: sensor.ActiveVersion,
		NewVersion:  version,
		RequestedAt: s.now(),
		CompletedAt: 0,
		Outcome:     "pending",
		Digest:      digest,
	}
	if err := s.store.PutRotation(rotation); err != nil {
		return model.RotationResult{}, err
	}
	return model.RotationResult{Rotation: rotation, Summary: model.SecretSummary{SensorID: req.SensorID, Version: version, Digest: digest, Length: len(req.Secret), Accepted: true}, Sensor: sensor}, s.log("rotation.prepared", req.SensorID, "success", "secret wrapped", map[string]string{"envelope": envelope.ID})
}

func (s *Service) CommitWrappedSecret(sensorID string, version int64, transport string, actor string) (model.SecretSummary, error) {
	if strings.TrimSpace(sensorID) == "" {
		return model.SecretSummary{}, errors.New("sensor id is required")
	}
	if version < 1 {
		return model.SecretSummary{}, errors.New("version is required")
	}
	if strings.TrimSpace(actor) == "" {
		return model.SecretSummary{}, errors.New("actor is required")
	}
	controllerID, decodedSensorID, decodedVersion, ciphertext, digest, err := crypto.DecodeTransport(transport)
	if err != nil {
		return model.SecretSummary{}, err
	}
	if controllerID != s.clockID || decodedSensorID != sensorID || decodedVersion != version {
		return model.SecretSummary{}, errors.New("transport envelope does not match request")
	}
	envelope, err := s.store.GetEnvelope(fmt.Sprintf("env-%s-%d", sensorID, version))
	if err != nil {
		return model.SecretSummary{}, err
	}
	secret, err := s.codec.Open(sensorID, version, ciphertext)
	if err != nil {
		return model.SecretSummary{}, err
	}
	verifiedDigest := crypto.DigestSecret(secret)
	if !model.CompareDigest(verifiedDigest, digest) {
		return model.SecretSummary{}, fmt.Errorf("sensor rejected envelope for %s", sensorID)
	}
	if !strings.EqualFold(envelope.Digest, digest) {
		return model.SecretSummary{}, errors.New("stored envelope digest no longer matches transport")
	}
	return s.applyRotation(sensorID, version, digest, len(secret), actor)
}

func (s *Service) EnvelopeTransport(sensorID string, version int64) (string, error) {
	envelope, err := s.store.GetEnvelope(fmt.Sprintf("env-%s-%d", sensorID, version))
	if err != nil {
		return "", err
	}
	if envelope.TransportTag == "" {
		return "", errors.New("envelope transport is empty")
	}
	return envelope.TransportTag, nil
}

func (s *Service) FinalizeRotation(sensorID string, version int64, digest string, length int, actor string) (model.SecretSummary, error) {
	if sensorID == "" || version < 1 || digest == "" || length < 1 {
		return model.SecretSummary{}, errors.New("rotation finalization fields are incomplete")
	}
	return s.applyRotation(sensorID, version, digest, length, actor)
}

func (s *Service) applyRotation(sensorID string, version int64, digest string, length int, actor string) (model.SecretSummary, error) {
	var summary model.SecretSummary
	sensor, err := s.store.GetSensor(sensorID)
	if err != nil {
		return model.SecretSummary{}, err
	}
	sensor.ActiveVersion = version
	sensor.Status = "active"
	sensor.UpdatedAt = s.now()
	if err := s.store.PutSensor(sensor); err != nil {
		return model.SecretSummary{}, err
	}
	rotations, err := s.store.ListRotations(sensorID)
	if err != nil {
		return model.SecretSummary{}, err
	}
	for index := range rotations {
		if rotations[index].NewVersion == version {
			rotations[index].Outcome = "success"
			rotations[index].CompletedAt = s.now()
			rotations[index].Digest = digest
			if err := s.store.PutRotation(rotations[index]); err != nil {
				return model.SecretSummary{}, err
			}
			break
		}
	}
	envelope, err := s.store.GetEnvelope(fmt.Sprintf("env-%s-%d", sensorID, version))
	if err == nil {
		envelope.State = "active"
		if err := s.store.PutEnvelope(envelope); err != nil {
			return model.SecretSummary{}, err
		}
	}
	summary = model.SecretSummary{SensorID: sensorID, Version: version, Digest: digest, Length: length, Accepted: true}
	if err := s.log("rotation.completed", sensorID, "success", "secret committed", map[string]string{"actor": actor, "version": fmt.Sprintf("%d", version)}); err != nil {
		return model.SecretSummary{}, err
	}
	return summary, nil
}

func (s *Service) InspectKeyLength(sensorID string) (int, error) {
	sensor, err := s.store.GetSensor(sensorID)
	if err != nil {
		return 0, err
	}
	if sensor.KeyLength < 1 {
		return 0, errors.New("sensor key length is unavailable")
	}
	return sensor.KeyLength, nil
}

func (s *Service) AuditRotations(sensorID string) ([]model.AuditEntry, error) {
	entries, err := s.store.ListAudits(sensorID, "")
	if err != nil {
		return nil, err
	}
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].At == entries[j].At {
			return entries[i].ID < entries[j].ID
		}
		return entries[i].At < entries[j].At
	})
	return entries, nil
}

func (s *Service) ListActiveSensors() ([]model.Sensor, error) {
	sensors, err := s.store.ListSensors()
	if err != nil {
		return nil, err
	}
	var result []model.Sensor
	for _, sensor := range sensors {
		if sensor.Status == "active" {
			result = append(result, sensor)
		}
	}
	return result, nil
}

func (s *Service) log(event, sensorID, outcome, message string, metadata map[string]string) error {
	if metadata == nil {
		metadata = map[string]string{}
	}
	if event == "" || sensorID == "" {
		return errors.New("audit event and sensor id are required")
	}
	entry := model.AuditEntry{
		ID:       fmt.Sprintf("audit-%s-%s-%d", sensorID, event, s.now()),
		Event:    event,
		SensorID: sensorID,
		Actor:    s.clockID,
		Outcome:  outcome,
		Message:  message,
		At:       s.now(),
		Metadata: metadata,
	}
	if err := s.store.PutAudit(entry); err != nil {
		return err
	}
	return s.logger.Record(entry)
}
