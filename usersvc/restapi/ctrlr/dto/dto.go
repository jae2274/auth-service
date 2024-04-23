package dto

import (
	"github.com/jae2274/goutils/enum"
)

type AuthenticateResponse struct {
	AuthToken string
}

type AuthCodeUrlsResponse struct {
	AuthCodeUrls []*AuthCodeUrlRes `json:"authCodeUrls"`
}

type AuthCodeUrlRes struct {
	AuthServer string `json:"authServer"`
	Url        string `json:"url"`
}

type SignInRequest struct {
	AuthToken string `json:"authToken"`
}

type SignInStatusValues struct{}
type SignInStatus enum.Enum[SignInStatusValues]

const (
	SignInSuccess          = SignInStatus("success")
	SignInNewUser          = SignInStatus("new_user")
	SignInRequireAgreement = SignInStatus("require_agreement")
)

func (SignInStatusValues) Values() []string {
	return []string{string(SignInSuccess), string(SignInNewUser)}
}

type SignInResponse struct {
	SignInStatus        SignInStatus         `json:"signInStatus"`
	SuccessRes          *SignInSuccessRes    `json:"successRes"`
	NewUserRes          *SignInNewUserRes    `json:"newUserRes"`
	RequireAgreementRes *RequireAgreementRes `json:"requireAgreementRes"`
}

type SignInSuccessRes struct {
	Username     string   `json:"username"`
	Roles        []string `json:"roles"`
	GrantType    string   `json:"grantType"`
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
}

type SignInNewUserRes struct {
	Email      string          `json:"email"`
	Agreements []*AgreementRes `json:"agreements"`
}

type RequireAgreementRes struct {
	Agreements []*AgreementRes `json:"agreements"`
}

type AgreementRes struct {
	AgreementId int    `json:"agreementId"`
	Required    bool   `json:"required"`
	Summary     string `json:"summary"`
	Priority    int    `json:"priority"`
}

type SignUpRequest struct {
	AuthToken  string              `json:"authToken"`
	Username   string              `json:"username"`
	Agreements []*UserAgreementReq `json:"agreements"`
}

type UserAgreementReq struct {
	AgreementId int  `json:"agreementId"`
	IsAgree     bool `json:"isAgree"`
}
