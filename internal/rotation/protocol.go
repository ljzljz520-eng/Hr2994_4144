package rotation

import (
	"encoding/json"
	"errors"
	"fmt"

	"industrial-key-rotation/internal/model"
)

func EncodeRotationResult(result model.RotationResult) ([]byte, error) {
	if err := result.Rotation.Validate(); err != nil {
		return nil, err
	}
	if err := result.Summary.Validate(); err != nil {
		return nil, err
	}
	return json.Marshal(result)
}

func DecodeRotationResult(data []byte) (model.RotationResult, error) {
	if len(data) == 0 {
		return model.RotationResult{}, errors.New("rotation result is empty")
	}
	var result model.RotationResult
	if err := json.Unmarshal(data, &result); err != nil {
		return result, fmt.Errorf("decode rotation result: %w", err)
	}
	if err := result.Rotation.Validate(); err != nil {
		return result, err
	}
	return result, result.Summary.Validate()
}

func ExplainOutcome(result model.RotationResult) string {
	if result.Summary.Accepted {
		return fmt.Sprintf("sensor %s active at version %d", result.Summary.SensorID, result.Summary.Version)
	}
	return fmt.Sprintf("sensor %s rejected version %d", result.Summary.SensorID, result.Summary.Version)
}
