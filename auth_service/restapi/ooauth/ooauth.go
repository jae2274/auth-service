package ooauth

import (
	"context"

	"github.com/jae2274/auth-service/auth_service/common/domain"

	"golang.org/x/oauth2"
)

type UserInfo struct {
	AuthorizedBy domain.AuthorizedBy `validate:"nonzero"`
	AuthorizedID string              `validate:"nonzero"`
	Email        string              `validate:"nonzero"`
	Username     string              `validate:"nonzero"`
}

type OauthToken struct {
	UserInfo *UserInfo     `json:"userInfo"`
	Token    *oauth2.Token `json:"authToken"`
}

type Ooauth interface {
	GetAuthServer() domain.AuthorizedBy
	GetLoginURL(state string) string
	Oauth2Config() *oauth2.Config
	GetToken(ctx context.Context, code string) (*oauth2.Token, error)
	GetUserInfo(ctx context.Context, authToken *OauthToken) (*UserInfo, error)
}
