// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package sqlc

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Transition struct {
	ID         uuid.UUID     `json:"id"`
	WorkflowID uuid.UUID     `json:"workflow_id"`
	FromNode   string        `json:"from_node"`
	ToNode     string        `json:"to_node"`
	Label      string        `json:"label"`
	CreatedAt  time.Time     `json:"created_at"`
	Previous   uuid.NullUUID `json:"previous"`
	Next       uuid.NullUUID `json:"next"`
}

type Workflow struct {
	ID           uuid.UUID       `json:"id"`
	CurrentNode  string          `json:"current_node"`
	Status       string          `json:"status"`
	Graph        json.RawMessage `json:"graph"`
	CreatedAt    time.Time       `json:"created_at"`
	NextActionAt time.Time       `json:"next_action_at"`
}
