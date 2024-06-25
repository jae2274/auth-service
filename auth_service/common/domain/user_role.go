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
	AuthorityName string     `json:"authorityName"`
	Summary       string     `json:"summary"`
	ExpiryDate    *time.Time `json:"expiryDate"`
}
