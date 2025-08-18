package data

import (
	"time"

	"github.com/google/uuid"
)

type Filter struct {
	ID             *uuid.UUID `json:"id"`
	Issuer         *string    `json:"issuer"`
	IssuedAtFrom   *time.Time `json:"issuedAtFrom"`
	IssuedAtTo     *time.Time `json:"issuedAtTo"`
	ExpirationFrom *time.Time `json:"expirationFrom"`
	ExpirationTo   *time.Time `json:"expirationTo"`

	OrderBy         []string  `json:"orderBy,omitzero"`
	OrderBySafeList []string  `json:"orderBySafeList,omitzero"`
	LastSeen        uuid.UUID `json:"lastSeen,omitzero"`
	PageSize        int       `json:"pageSize,omitzero"`
}
