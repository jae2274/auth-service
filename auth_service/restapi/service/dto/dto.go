package dto

import (
	"time"

	"github.com/jae2274/goutils/enum"
)

type AuthStatusValues struct{}
type AuthStatus enum.Enum[AuthStatusValues]

const (
	AuthSuccess = AuthStatus("success")
	AuthFailed  = AuthStatus("failed")
)

func (AuthStatusValues) Values() []string {
	return []string{string(AuthSuccess), string(AuthFailed)}
}

type AfterAuthViewVars struct {
	AuthStatus AuthStatus
	AuthToken  string
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
	SignInSuccess = SignInStatus("success")
	SignInFailed  = SignInStatus("failed")
	SignInNewUser = SignInStatus("new_user")
)

func (SignInStatusValues) Values() []string {
	return []string{string(SignInSuccess), string(SignInFailed), string(SignInNewUser)}
}

type SignInResponse struct {
	SignInStatus SignInStatus      `json:"signInStatus"`
	SuccessRes   *SignInSuccessRes `json:"successRes"`
	NewUserRes   *SignInNewUserRes `json:"newUserRes"`
}

type SignInSuccessRes struct {
	GrantType    string `json:"grantType"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type SignInNewUserRes struct {
	Email      string          `json:"email"`
	Agreements []*AgreementRes `json:"agreements"`
}

type AgreementRes struct {
	AgreementCode string `json:"agreementCode"`
	IsRequired    bool   `json:"isRequired"`
	Summary       string `json:"summary"`
	Priority      int    `json:"priority"`
}

type SignUpRequest struct {
	AuthToken  string              `json:"authToken"`
	Agreements []*UserAgreementReq `json:"agreements"`
}

type UserAgreementReq struct {
	AgreementID int  `json:"agreementID"`
	IsAgree     bool `json:"isAgree"`
}

type UserRoleReq struct {
	RoleName   string    `json:"roleName"`
	ExpiryDate time.Time `json:"expiryDate"`
}
