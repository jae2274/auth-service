package domain

import (
	"github.com/jae2274/auth-service/auth_service/models"
	"github.com/jae2274/goutils/enum"
)

type AuthorizedByValues struct{}
type AuthorizedBy = enum.Enum[AuthorizedByValues]

const (
	GOOGLE              = AuthorizedBy(models.UserAuthorizedByGOOGLE)
	AuthorizedByDELETED = AuthorizedBy(models.UserAuthorizedByDELETED)
)

func (AuthorizedByValues) Values() []string {
	return []string{string(GOOGLE)}
}

type StatusValues struct{}
type Status = enum.Enum[StatusValues]

const (
	ACTIVE    = Status(models.UserStatusACTIVE)
	SUSPENDED = Status(models.UserStatusSUSPENDED)
	DELETED   = Status(models.UserStatusDELETED)
)

func (StatusValues) Values() []string {
	return []string{string(ACTIVE), string(SUSPENDED), string(DELETED)}
}

type User struct {
	UserID           int              `json:"userId"`
	AuthorizedBy     AuthorizedBy     `json:"authorizedBy"`
	AuthorizedID     string           `json:"authorizedId"`
	UserName         string           `json:"userName"`
	Email            string           `json:"email"`
	Authorities      []*UserAuthority `json:"authorities"`
	CreatedUnixMilli int64            `json:"createdUnixMilli"`
}
