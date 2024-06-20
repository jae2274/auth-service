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
	AuthorityName string     `json:"authorityName"`
	ExpiryDate    *time.Time `json:"expiryDate"`
}
