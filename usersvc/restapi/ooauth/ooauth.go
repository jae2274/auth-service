package ooauth

import (
	"context"
	"userService/usersvc/common/domain"

	"golang.org/x/oauth2"
)

type UserInfo struct {
	AuthorizedBy domain.AuthorizedBy `validate:"nonzero"`
	AuthorizedID string              `validate:"nonzero"`
	Email        string              `validate:"nonzero"`
}

type Ooauth interface {
	GetLoginURL(state string) string
	Oauth2Config() *oauth2.Config
	GetToken(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, authToken *oauth2.Token) (*UserInfo, error)
}
