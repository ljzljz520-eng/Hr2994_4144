package persistence

import (
	"fmt"

	"go.etcd.io/bbolt"
	"industrial-key-rotation/internal/model"
)

func (s *Store) PutSensor(value model.Sensor) error {
	data, err := model.EncodeSensor(value)
	if err != nil {
		return err
	}
	return s.put(sensorsBucket, value.ID, data)
}

func (s *Store) GetSensor(id string) (model.Sensor, error) {
	data, err := s.get(sensorsBucket, id)
	if err != nil {
		return model.Sensor{}, err
	}
	return model.DecodeSensor(data)
}

func (s *Store) PutEnvelope(value model.KeyEnvelope) error {
	data, err := model.EncodeEnvelope(value)
	if err != nil {
		return err
	}
	return s.put(envelopesBucket, value.ID, data)
}

func (s *Store) GetEnvelope(id string) (model.KeyEnvelope, error) {
	data, err := s.get(envelopesBucket, id)
	if err != nil {
		return model.KeyEnvelope{}, err
	}
	return model.DecodeEnvelope(data)
}

func (s *Store) PutRotation(value model.Rotation) error {
	data, err := model.EncodeRotation(value)
	if err != nil {
		return err
	}
	return s.put(rotationsBucket, value.ID, data)
}

func (s *Store) GetRotation(id string) (model.Rotation, error) {
	data, err := s.get(rotationsBucket, id)
	if err != nil {
		return model.Rotation{}, err
	}
	return model.DecodeRotation(data)
}

func (s *Store) PutAudit(value model.AuditEntry) error {
	data, err := model.EncodeAudit(value)
	if err != nil {
		return err
	}
	return s.put(auditBucket, value.ID, data)
}

func (s *Store) GetAudit(id string) (model.AuditEntry, error) {
	data, err := s.get(auditBucket, id)
	if err != nil {
		return model.AuditEntry{}, err
	}
	return model.DecodeAudit(data)
}

func (s *Store) put(bucket []byte, id string, data []byte) error {
	if id == "" {
		return fmt.Errorf("record id is required")
	}
	return s.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("bucket %s is missing", bucket)
		}
		return b.Put([]byte(id), data)
	})
}

func (s *Store) get(bucket []byte, id string) ([]byte, error) {
	if id == "" {
		return nil, fmt.Errorf("record id is required")
	}
	var result []byte
	err := s.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("bucket %s is missing", bucket)
		}
		value := b.Get([]byte(id))
		if value == nil {
			return ErrNotFound
		}
		result = append([]byte(nil), value...)
		return nil
	})
	return result, err
}
