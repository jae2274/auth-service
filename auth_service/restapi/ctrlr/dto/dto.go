package dto

import (
	"time"

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

type UserInfoRequest struct {
	AuthToken string `json:"authToken"`
}

type UserInfoResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type SignInRequest struct {
	AuthToken            string              `json:"authToken"`
	AdditionalAgreements []*UserAgreementReq `json:"additionalAgreements"`
}

type SignInStatusValues struct{}
type SignInStatus enum.Enum[SignInStatusValues]

const (
	SignInSuccess             = SignInStatus("success")
	SignInNewUser             = SignInStatus("new_user")
	SignInNecessaryAgreements = SignInStatus("necessary_agreements")
)

func (SignInStatusValues) Values() []string {
	return []string{string(SignInSuccess), string(SignInNewUser)}
}

type SignInResponse struct {
	SignInStatus           SignInStatus                  `json:"signInStatus"`
	SuccessRes             *SignInSuccessRes             `json:"successRes"`
	NewUserRes             *SignInNewUserRes             `json:"newUserRes"`
	NecessaryAgreementsRes *SignInNecessaryAgreementsRes `json:"necessaryAgreementsRes"`
}

type SignInSuccessRes struct {
	Username     string   `json:"username"`
	Authorities  []string `json:"authorities"`
	GrantType    string   `json:"grantType"`
	AccessToken  string   `json:"accessToken"`
	RefreshToken string   `json:"refreshToken"`
}

type SignInNewUserRes struct {
	Email      string          `json:"email"`
	Username   string          `json:"username"`
	Agreements []*AgreementRes `json:"agreements"`
}

type SignInNecessaryAgreementsRes struct {
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

type RefreshJwtRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type RefreshJwtResponse struct {
	AccessToken string   `json:"accessToken"`
	Authorities []string `json:"authorities"`
}

type AddAuthorityRequest struct {
	UserId           int                 `json:"userId"`
	AddedAuthorities []*UserAuthorityReq `json:"authorities"`
}

type UserAuthorityReq struct {
	AuthorityID    int       `json:"-"`
	AuthorityName  string    `json:"authorityName"`
	ExpiryDuration *Duration `json:"expiryDate"`
}

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Duration(d).String() + `"`), nil
}

func (d *Duration) UnmarshalJSON(data []byte) error {
	duration, err := time.ParseDuration(string(data[1 : len(data)-1]))
	if err != nil {
		return err
	}

	*d = Duration(duration)

	return nil
}