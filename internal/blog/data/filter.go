package data

import (
	"time"

	"github.com/google/uuid"
)

type Filter struct {
	ID            *uuid.UUID `json:"id"`
	Title         *string    `json:"title"`
	Content       *string    `json:"content"`
	Published     *bool      `json:"published"`
	CreatedAtFrom *time.Time `json:"createdAtFrom"`
	CreatedAtTo   *time.Time `json:"createdAtTo"`
	UpdatedAtFrom *time.Time `json:"updatedAtFrom"`
	UpdatedAtTo   *time.Time `json:"updatedAtTo"`
	Deleted       *bool      `json:"deleted"`
	DeletedAtFrom *bool      `json:"deletedAtFrom"`
	DeletedAtTo   *bool      `json:"deletedAtTo"`

	OrderBy         []string  `json:"orderBy,omitzero"`
	OrderBySafeList []string  `json:"orderBySafeList,omitzero"`
	LastSeen        uuid.UUID `json:"lastSeen,omitzero"`
	PageSize        int       `json:"pageSize,omitzero"`
}
