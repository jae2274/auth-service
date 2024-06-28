package dto

import (
	"github.com/jae2274/auth-service/auth_service/common/domain"
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
	UserId      int                 `json:"userId"`
	Authorities []*UserAuthorityReq `json:"authorities"`
}

type UserAuthorityReq struct {
	AuthorityID      int    `json:"-"`
	AuthorityCode    string `json:"authorityCode"`
	ExpiryDurationMS *int64 `json:"expiryDurationMS"`
}

type RemoveAuthorityRequest struct {
	UserId        int    `json:"userId"`
	AuthorityCode string `json:"authorityCode"`
}

type TicketStatusValues struct{}
type TicketStatus enum.Enum[TicketStatusValues]

const (
	NOT_EXISTED       = TicketStatus("not_existed")
	ALREADY_USED      = TicketStatus("already_used")
	SUCCESSFULLY_USED = TicketStatus("successfully_used")
)

func (TicketStatusValues) Values() []string {
	return []string{string(NOT_EXISTED), string(ALREADY_USED), string(SUCCESSFULLY_USED)}
}

type UseTicketResponse struct {
	TicketStatus TicketStatus `json:"ticketStatus"`
	AccessToken  *string      `json:"accessToken"`

	Authorities []string `json:"authorities"`
	// UserAuthorities []*domain.UserAuthority `json:"appliedAuthorities"`
}

type GetAllAuthoritiesResponse struct {
	Authorities []*AuthorityRes `json:"authorities"`
}
type AuthorityRes struct {
	AuthorityCode string `json:"authorityCode"`
	AuthorityName string `json:"authorityName"`
	Summary       string `json:"summary"`
}

type GetAllTicketsResponse struct {
	Tickets []*Ticket `json:"tickets"`
}

type GetAllUsersResponse struct {
	Users []*domain.User `json:"users"`
}
