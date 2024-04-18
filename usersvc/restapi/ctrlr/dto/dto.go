package dto

import "github.com/jae2274/goutils/enum"

type LoginStatusValue struct{}
type LoginStatus enum.Enum[LoginStatusValue]

const (
	LoginSuccess = LoginStatus("success")
	LoginFailed  = LoginStatus("failed")
	LoginNewUser = LoginStatus("new_user")
)

func (LoginStatusValue) Values() []string {
	return []string{string(LoginSuccess), string(LoginFailed), string(LoginNewUser)}
}

type AfterLoginViewVars struct {
	LoginStatus LoginStatus

	GrantType    string
	AccessToken  string
	RefreshToken string

	AuthToken string
	Email     string
}
