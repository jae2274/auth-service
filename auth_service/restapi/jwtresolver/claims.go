package jwtresolver

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/jae2274/auth-service/auth_service/common/domain"
)

type CustomClaims struct {
	UserId       string `validate:"nonzero"`
	AuthorizedBy domain.AuthorizedBy
	AuthorizedID string

	Authorities []string
	jwt.RegisteredClaims
}

func (c *CustomClaims) HasRole(role string) bool {
	for _, r := range c.Authorities {
		if r == role {
			return true
		}
	}
	return false
}
