package entity

import (
	"time"
	"userService/usersvc/common/domain"

	"github.com/jae2274/goutils/enum"
)

type UserVO struct {
	UserID       int64
	AuthorizedBy domain.AuthorizedBy
	AuthorizedID string
	Email        string
	CreateDate   time.Time
}

type GrantedTypeValues struct{}
type GrantedType = enum.Enum[GrantedTypeValues]

const (
	ADMIN  = GrantedType("ADMIN")
	TICKET = GrantedType("TICKET")
)

func (GrantedTypeValues) Values() []string {
	return []string{string(ADMIN), string(TICKET)}
}

type UserRoleVO struct {
	UserID      int64
	RoleName    string
	GrantedType GrantedType
	GrantedBy   int64
	ExpiryDate  *time.Time
}
