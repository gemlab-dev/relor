package kv

import (
	"encoding/binary"
	"fmt"

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
