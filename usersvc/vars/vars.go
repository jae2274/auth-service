package vars

import (
	"fmt"
	"os"
	"strconv"
)

type DBVars struct {
	Host     string
	Port     int64
	Name     string
	Username string
	Password string
}

type Vars struct {
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectUrl  string
	SecretKey          string
	DbVars             *DBVars
	ApiPort            int
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

	googleRedirectUrl, err := getFromEnv("GOOGLE_REDIRECT_URL")
	if err != nil {
		return nil, NotExistedVar("GOOGLE_REDIRECT_URL")
	}

	SecretKey, err := getFromEnv("SECRET_KEY")
	if err != nil {
		return nil, NotExistedVar("SECRET_KEY")
	}

	apiPort, err := getFromEnv("API_PORT")
	if err != nil {
		return nil, err
	}

	userSvcPortInt, err := strconv.ParseInt(apiPort, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("REST_API_PORT is not integer.\tREST_API_PORT: %s", apiPort)
	}

	dbVars, err := getDBVars()
	if err != nil {
		return nil, err
	}

	return &Vars{
		GoogleClientID:     googleClientID,
		GoogleClientSecret: googleClientSecret,
		GoogleRedirectUrl:  googleRedirectUrl,
		SecretKey:          SecretKey,
		ApiPort:            int(userSvcPortInt),
		DbVars:             dbVars,
	}, nil
}

func getDBVars() (*DBVars, error) {
	dbHost, err := getFromEnv("DB_HOST")
	if err != nil {
		return nil, err
	}

	dbPortStr, err := getFromEnv("DB_PORT")
	if err != nil {
		return nil, err
	}

	dbPort, err := strconv.ParseInt(dbPortStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("DB_PORT is not integer.\tDB_PORT: %s", dbPortStr)
	}

	dbName, err := getFromEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	dbUsername, err := getFromEnv("DB_USERNAME")
	if err != nil {
		return nil, err
	}

	dbPassword, err := getFromEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	return &DBVars{
		Host:     dbHost,
		Port:     dbPort,
		Name:     dbName,
		Username: dbUsername,
		Password: dbPassword,
	}, nil
}

func getFromEnv(envVar string) (string, error) {
	ev := os.Getenv(envVar)

	if ev == "" {
		return "", fmt.Errorf("%s is not existed", envVar)
	}

	return ev, nil
}
