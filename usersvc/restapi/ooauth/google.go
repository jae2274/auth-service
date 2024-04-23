package ooauth

import (
	"context"
	"encoding/json"
	"io"
	"userService/usersvc/common/domain"

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

func (g *GoogleOauth) GetAuthServer() domain.AuthorizedBy {
	return domain.GOOGLE
}

func (g *GoogleOauth) GetLoginURL(state string) string {
	return g.Oauth2Config().AuthCodeURL(state)
}

func (g *GoogleOauth) Oauth2Config() *oauth2.Config {
	return g.oauthConfig
}

func (g *GoogleOauth) GetToken(ctx context.Context, code string) (*OauthToken, error) {
	token, err := g.Oauth2Config().Exchange(ctx, code)
	if err != nil {
		return nil, err
	}
	return &OauthToken{
		AuthServer: g.GetAuthServer(),
		Token:      token,
	}, nil
}

type user struct {
	Sub   string `json:"sub"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (g *GoogleOauth) GetUserInfo(ctx context.Context, authToken *OauthToken) (*UserInfo, error) {
	client := g.Oauth2Config().Client(ctx, authToken.Token)
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
		AuthorizedBy: domain.GOOGLE,
		AuthorizedID: authUser.Sub,
		Email:        authUser.Email,
		Username:     authUser.Name,
	}, nil
}
