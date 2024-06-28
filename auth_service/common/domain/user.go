package domain

import (
	"github.com/jae2274/goutils/enum"
)

type AuthorizedByValues struct{}
type AuthorizedBy = enum.Enum[AuthorizedByValues]

const (
	GOOGLE = AuthorizedBy("GOOGLE")
)

func (AuthorizedByValues) Values() []string {
	return []string{string(GOOGLE)}
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
