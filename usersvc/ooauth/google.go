package ooauth

import (
	"context"
	"encoding/json"
	"io"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	scopeEmail   = "https://www.googleapis.com/auth/userinfo.email"
	scopeProfile = "https://www.googleapis.com/auth/userinfo.profile"

	userInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"
)

type GoogleOauth struct {
	oauthConfig *oauth2.Config
}

func NewGoogleOauth(clientId, clientSecret, redirectUrl string) Ooauth {
	return &GoogleOauth{
		oauthConfig: &oauth2.Config{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			RedirectURL:  redirectUrl,
			Scopes:       []string{scopeEmail, scopeProfile},
			Endpoint:     google.Endpoint,
		},
	}
}

func (g *GoogleOauth) GetLoginURL(state string) string {
	return g.Oauth2Config().AuthCodeURL(state)
}

func (g *GoogleOauth) Oauth2Config() *oauth2.Config {
	return g.oauthConfig
}

func (g *GoogleOauth) GetToken(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := g.Oauth2Config().Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return token, nil
}

type user struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (g *GoogleOauth) GetUserInfo(ctx context.Context, authToken *oauth2.Token) (*UserInfo, error) {
	client := g.Oauth2Config().Client(ctx, authToken)
	userInfoResp, err := client.Get(userInfoURL)
	if err != nil {
		return nil, err
	}
	defer userInfoResp.Body.Close()
	userInfo, err := io.ReadAll(userInfoResp.Body)
	if err != nil {
		return nil, err
	}
	var authUser user
	json.Unmarshal(userInfo, &authUser)

	return &UserInfo{
		Email: authUser.Email,
	}, nil
}
