package vars

import (
	"fmt"
	"os"
	"strconv"
)

type Vars struct {
	GoogleClientID     string
	GoogleClientSecret string
	RedirectURL        string
	Port               int
}

type ErrNotExistedVar struct {
	VarName string
}

func NotExistedVar(varName string) *ErrNotExistedVar {
	return &ErrNotExistedVar{VarName: varName}
}

func (e *ErrNotExistedVar) Error() string {
	return fmt.Sprintf("%s is not existed", e.VarName)
}

func Variables() (*Vars, error) {

	googleClientID, err := getFromEnv("GOOGLE_CLIENT_ID")
	if err != nil {
		return nil, NotExistedVar("GOOGLE_CLIENT_ID")
	}

	googleClientSecret, err := getFromEnv("GOOGLE_CLIENT_SECRET")
	if err != nil {
		return nil, NotExistedVar("GOOGLE_CLIENT_SECRET")
	}

	redirectURL, err := getFromEnv("REDIRECT_URL")
	if err != nil {
		return nil, NotExistedVar("REDIRECT_URL")
	}

	userSvcPort, err := getFromEnv("USER_SVC_PORT")
	if err != nil {
		return nil, err
	}

	userSvcPortInt, err := strconv.ParseInt(userSvcPort, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("REST_API_PORT is not integer.\tREST_API_PORT: %s", userSvcPort)
	}

	return &Vars{
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		RedirectURL:        redirectURL,
		Port:               int(userSvcPortInt),
	}, nil
}

func getFromEnv(envVar string) (string, error) {
	ev := os.Getenv(envVar)

	if ev == "" {
		return "", fmt.Errorf("%s is not existed", envVar)
	}

	return ev, nil
}
