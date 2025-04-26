package storage

import "github.com/google/uuid"

type NextAction struct {
	ID                  uuid.UUID
	Label               string
	CurrentTransition   uuid.UUID
	LastKnownTransition uint64
}
