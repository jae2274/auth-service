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

type OauthToken struct {
	AuthServer domain.AuthorizedBy `json:"authServer"`
	Token      *oauth2.Token       `json:"authToken"`
}

type Ooauth interface {
	GetAuthServer() domain.AuthorizedBy
	GetLoginURL(state string) string
	Oauth2Config() *oauth2.Config
	GetToken(ctx context.Context, code string) (*OauthToken, error)
	GetUserInfo(ctx context.Context, authToken *OauthToken) (*UserInfo, error)
}
