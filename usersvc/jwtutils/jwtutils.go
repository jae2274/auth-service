package jwtutils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserId string
	Roles  []string
	jwt.RegisteredClaims
}

type JwtResolver struct {
	secretKey []byte
}

func NewJwtUtils(secretKey string) *JwtResolver {
	return &JwtResolver{
		secretKey: []byte(secretKey),
	}
}

func (j *JwtResolver) CreateToken(userId string, email string, roles []string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&CustomClaims{
			UserId: userId,
			Roles:  roles,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "careerhub.jyo-liar.com", //TODO: 임의 설정
				Subject:   email,
				Audience:  []string{"careerhub.jyo-liar.com"},          //TODO: 임의 설정
				ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)), //TODO: 임의 설정
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
			},
		},
	)

	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
