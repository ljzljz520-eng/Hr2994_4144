package policy

import (
	"errors"
	"fmt"
	"sort"

	"industrial-key-rotation/internal/model"
)

type Lifecycle struct {
	CurrentVersion int64
	PendingVersion int64
	State          string
}

func BuildLifecycle(sensor model.Sensor, rotations []model.Rotation) Lifecycle {
	lifecycle := Lifecycle{CurrentVersion: sensor.ActiveVersion, State: sensor.Status}
	for _, rotation := range rotations {
		if rotation.NewVersion > lifecycle.PendingVersion && rotation.Outcome == "pending" {
			lifecycle.PendingVersion = rotation.NewVersion
		}
	}
	return lifecycle
}

func (l Lifecycle) Validate() error {
	if l.CurrentVersion < 0 || l.PendingVersion < 0 {
		return errors.New("lifecycle versions cannot be negative")
	}
	if l.PendingVersion > 0 && l.PendingVersion <= l.CurrentVersion {
		return errors.New("pending version must be newer than current")
	}
	if l.State == "" {
		return errors.New("lifecycle state is required")
	}
	return nil
}

func (l Lifecycle) NextVersion() int64 {
	if l.PendingVersion > l.CurrentVersion {
		return l.PendingVersion + 1
	}
	return l.CurrentVersion + 1
}

func (l Lifecycle) CanActivate(version int64) error {
	if version < 1 {
		return errors.New("activation version must be positive")
	}
	if version <= l.CurrentVersion {
		return fmt.Errorf("version %d is not newer than %d", version, l.CurrentVersion)
	}
	if l.State == "suspended" {
		return errors.New("suspended lifecycle cannot activate")
	}
	return nil
}

func SortRotations(rotations []model.Rotation) []model.Rotation {
	result := append([]model.Rotation(nil), rotations...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].NewVersion == result[j].NewVersion {
			return result[i].ID < result[j].ID
		}
		return result[i].NewVersion < result[j].NewVersion
	})
	return result
}

func PendingRotations(rotations []model.Rotation) []model.Rotation {
	var result []model.Rotation
	for _, rotation := range rotations {
		if rotation.Outcome == "pending" {
			result = append(result, rotation)
		}
	}
	return SortRotations(result)
}
