package persistence

import (
	"sort"

	"go.etcd.io/bbolt"
	"industrial-key-rotation/internal/model"
)

func (s *Store) ListSensors() ([]model.Sensor, error) {
	var result []model.Sensor
	err := s.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(sensorsBucket).ForEach(func(_, value []byte) error {
			item, err := model.DecodeSensor(value)
			if err != nil {
				return err
			}
			result = append(result, item)
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}

func (s *Store) ListRotations(sensorID string) ([]model.Rotation, error) {
	var result []model.Rotation
	err := s.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(rotationsBucket).ForEach(func(_, value []byte) error {
			item, err := model.DecodeRotation(value)
			if err != nil {
				return err
			}
			if sensorID == "" || item.SensorID == sensorID {
				result = append(result, item)
			}
			return nil
		})
	})
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].RequestedAt == result[j].RequestedAt {
			return result[i].ID < result[j].ID
		}
		return result[i].RequestedAt < result[j].RequestedAt
	})
	return result, err
}

func (s *Store) ListAudits(sensorID string, event string) ([]model.AuditEntry, error) {
	var result []model.AuditEntry
	err := s.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(auditBucket).ForEach(func(_, value []byte) error {
			item, err := model.DecodeAudit(value)
			if err != nil {
				return err
			}
			if (sensorID == "" || item.SensorID == sensorID) && (event == "" || item.Event == event) {
				result = append(result, item)
			}
			return nil
		})
	})
	sort.SliceStable(result, func(i, j int) bool { return result[i].At < result[j].At })
	return result, err
}

func (s *Store) Count(kind string) (int, error) {
	bucket := bucketFor(kind)
	if bucket == nil {
		return 0, nil
	}
	count := 0
	err := s.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return nil
		}
		count = b.Stats().KeyN
		return nil
	})
	return count, err
}

func (s *Store) SnapshotCounts() (map[string]int, error) {
	result := make(map[string]int)
	for kind := range allBuckets() {
		count, err := s.Count(kind)
		if err != nil {
			return nil, err
		}
		result[kind] = count
	}
	return result, nil
}
