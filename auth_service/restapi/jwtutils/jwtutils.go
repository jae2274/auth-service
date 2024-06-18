package jwtutils

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jae2274/goutils/terr"
	"gopkg.in/validator.v2"
)

type CustomClaims struct {
	UserId string `validate:"nonzero"`
	Roles  []string
	jwt.RegisteredClaims
}

type JwtResolver struct {
	secretKey            []byte
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

type TokenInfo struct {
	GrantType    string `json:"grantType"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func NewJwtUtils(secretKey []byte) *JwtResolver {
	return &JwtResolver{
		secretKey:            secretKey,
		accessTokenDuration:  10 * time.Minute,
		refreshTokenDuration: 24 * time.Hour,
	}
}

func (j *JwtResolver) SetAccessTokenDuration(duration time.Duration) error {
	if duration < 0 {
		return terr.New("duration must be positive")
	}
	j.accessTokenDuration = duration
	return nil
}

func (j *JwtResolver) SetRefreshTokenDuration(duration time.Duration) error {
	if duration < 0 {
		return terr.New("duration must be positive")
	}
	j.refreshTokenDuration = duration
	return nil
}

func (j *JwtResolver) GetAccessTokenDuration() time.Duration {
	return j.accessTokenDuration
}

func (j *JwtResolver) GetRefreshTokenDuration() time.Duration {
	return j.refreshTokenDuration
}

func (j *JwtResolver) CreateToken(userId string, roles []string, createdAt time.Time) (*TokenInfo, error) {
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&CustomClaims{
			UserId: userId,
			Roles:  roles,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(createdAt.Add(j.accessTokenDuration)),
			},
		},
	).SignedString(j.secretKey)

	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&CustomClaims{
			UserId: userId,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(createdAt.Add(j.refreshTokenDuration)),
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

func (j *JwtResolver) ParseToken(tokenString string) (*CustomClaims, bool, error) {

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	jwtToken, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return &CustomClaims{}, false, nil
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return &CustomClaims{}, false, terr.New("invalid token. token is malformed")
	} else if err != nil {
		return &CustomClaims{}, false, terr.Wrap(err)
	} else if jwtToken.Valid {
		if claims, ok := jwtToken.Claims.(*CustomClaims); ok {
			if err := validator.Validate(claims); err != nil {
				return claims, false, err
			} else {
				return claims, true, nil
			}
		} else {
			return &CustomClaims{}, false, terr.New("invalid token. claims is not CustomClaims type")
		}
	} else {
		return &CustomClaims{}, false, nil
	}
}
