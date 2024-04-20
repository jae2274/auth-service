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

type TokenInfo struct {
	GrantType    string `json:"grantType"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func NewJwtUtils(secretKey []byte) *JwtResolver {
	return &JwtResolver{
		secretKey: secretKey,
	}
}

func (j *JwtResolver) CreateToken(userId string, userEmail string, roles []string) (*TokenInfo, error) {
	now := time.Now()

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&CustomClaims{
			UserId: userId,
			Roles:  roles,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "careerhub.jyo-liar.com", //TODO: 임의 설정
				Subject:   userEmail,
				Audience:  []string{"careerhub.jyo-liar.com"},          //TODO: 임의 설정
				ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Hour)), //TODO: 임의 설정
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
			},
		},
	).SignedString(j.secretKey)

	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&CustomClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(now.Add(24 * time.Hour)), //TODO: 임의 설정
			},
		},
	).SignedString(j.secretKey)

	if err != nil {
		return nil, err
	}

	return &TokenInfo{
		GrantType:    "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
