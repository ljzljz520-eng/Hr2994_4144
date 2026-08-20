package persistence

var (
	sensorsBucket   = []byte("sensors")
	envelopesBucket = []byte("envelopes")
	rotationsBucket = []byte("rotations")
	auditBucket     = []byte("audit_entries")
	metaBucket      = []byte("metadata")
)

func bucketNames() [][]byte {
	return [][]byte{sensorsBucket, envelopesBucket, rotationsBucket, auditBucket, metaBucket}
}

func bucketFor(kind string) []byte {
	switch kind {
	case "sensor":
		return sensorsBucket
	case "envelope":
		return envelopesBucket
	case "rotation":
		return rotationsBucket
	case "audit":
		return auditBucket
	default:
		return nil
	}
}

func allBuckets() map[string][]byte {
	return map[string][]byte{
		"sensor":   sensorsBucket,
		"envelope": envelopesBucket,
		"rotation": rotationsBucket,
		"audit":    auditBucket,
	}
}
