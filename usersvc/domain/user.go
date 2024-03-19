package domain

import (
	"time"

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
	UserID       string
	AuthorizedBy AuthorizedBy
	AuthorizedID string
	Email        string
	Roles        []string
	CreateDate   time.Time
}
