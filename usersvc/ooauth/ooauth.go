package ooauth

import (
	"context"
	"userService/usersvc/domain"

	"golang.org/x/oauth2"
)

type UserInfo struct {
	AuthorizedBy domain.AuthorizedBy
	AuthorizedID string
	Email        string
}

type Ooauth interface {
	GetLoginURL(state string) string
	Oauth2Config() *oauth2.Config
	GetToken(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, authToken *oauth2.Token) (*UserInfo, error)
}
