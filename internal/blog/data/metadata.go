package data

import "github.com/google/uuid"

type Metadata struct {
	LastSeen       uuid.UUID `json:"lastSeen,omitzero"`
	Next           bool      `json:"next"`
	ResponseLength int       `json:"responseLength"`
}
