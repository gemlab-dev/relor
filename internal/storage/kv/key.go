package kv

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TransitionKey generates a key for a transition using the given prefix and sequence ID.
// The sequence ID is negated to enforce descending order when sorted.
func TransitionKey(prefix []byte, seqID uint64) []byte {
	descBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(descBytes, ^seqID)
	return append(prefix, descBytes...)
}

// WorkflowKey generates a key for a workflow using the given prefix and workflow ID.
func WorkflowKey(prefix []byte, id uuid.UUID) ([]byte, error) {
	b, err := id.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow ID to bytes: %w", err)
	}
	return append(prefix, b...), nil
}

func SchedulePrefix(prefix []byte, when time.Time) ([]byte, error) {
	if when.IsZero() {
		return nil, fmt.Errorf("when is zero")
	}
	bts := make([]byte, 8)
	binary.BigEndian.PutUint64(bts, uint64(when.UnixNano()))
	return append(prefix, bts...), nil
}

func ScheduleKey(prefix []byte, when time.Time, id uuid.UUID) ([]byte, error) {
	pts, err := SchedulePrefix(prefix, when)
	if err != nil {
		return nil, fmt.Errorf("failed to create schedule key prefix: %w", err)
	}

	bid, err := id.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow ID to bytes: %w", err)
	}

	return append(pts, bid...), nil
}
