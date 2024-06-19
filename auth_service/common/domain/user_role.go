package domain

import (
	"time"

	"github.com/jae2274/goutils/enum"
)

type GrantedTypeValues struct{}

type GrantedType = enum.Enum[GrantedTypeValues]

const (
	ADMIN  = GrantedType("ADMIN")
	TICKET = GrantedType("TICKET")
)

func (GrantedTypeValues) Values() []string {
	return []string{string(ADMIN), string(TICKET)}
}

type UserAuthority struct {
	UserID        int        `json:"-"`
	AuthorityID   int        `json:"-"`
	AuthorityName string     `json:"authorityName"`
	ExpiryDate    *time.Time `json:"expiryDate"`
}
