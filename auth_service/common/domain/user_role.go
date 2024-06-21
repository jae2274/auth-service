package domain

import (
	"time"
)

const (
	AuthorityAdmin = "AUTHORITY_ADMIN"
)

type UserAuthority struct {
	UserID        int        `json:"-"`
	AuthorityID   int        `json:"-"`
	AuthorityCode string     `json:"authorityCode"`
	ExpiryDate    *time.Time `json:"expiryDate"`
}
