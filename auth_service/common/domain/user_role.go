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

type UserRole struct {
	RoleName       string         `json:"roleName"`
	ExpiryDuration *time.Duration `json:"expiryDate"`
}
