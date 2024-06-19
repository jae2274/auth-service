package jwtresolver

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserId      string `validate:"nonzero"`
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
