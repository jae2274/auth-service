package controller

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jae2274/auth-service/auth_service/common/domain"
	"github.com/jae2274/auth-service/auth_service/models"
	"github.com/jae2274/auth-service/auth_service/restapi/jwtresolver"
	"github.com/jae2274/auth-service/auth_service/restapi/ooauth"
	"github.com/jae2274/auth-service/auth_service/restapi/service"
	"github.com/jae2274/auth-service/test/tinit"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func newAuthorities() []*models.Authority {
	return []*models.Authority{
		{AuthorityName: domain.AuthorityAdmin, Summary: "관리자 권한"},
		{AuthorityName: "AUTHORITY_USER", Summary: "사용자 권한"},
		{AuthorityName: "AUTHORITY_GUEST", Summary: "게스트 권한"},
	}
}

func initAuthority(ctx context.Context, t *testing.T, db *sql.DB) []*models.Authority {
	authorities := newAuthorities()
	for _, authority := range authorities {
		err := authority.Insert(ctx, db, boil.Infer())
		require.NoError(t, err)
	}

	return authorities
}

func TestAdminController(t *testing.T) {
	ctx := context.Background()

	cancelFunc := tinit.RunTestApp(t)
	defer cancelFunc()

	rootUrl := initRootUrl(t)
	jwtResolver := initJwtResolver(t)
	userService := service.NewUserService(tinit.DB(t))

	signUpTestUser := func(t *testing.T) *models.User {
		user, err := userService.SignUp(ctx, &ooauth.UserInfo{AuthorizedBy: domain.GOOGLE, AuthorizedID: "authId", Email: "targetUser@test.com", Username: "target"}, nil)
		require.NoError(t, err)

		return user
	}
	sampleJsonBody := `
	{
		"userId": %d,
		"authorities": [
		  {
			"authorityName": "AUTHORITY_USER",
			"expiryDate": "720h"
		  },
		  {
			"authorityName": "AUTHORITY_GUEST"
		  }
		]
	  }
	`

	t.Run("return 401 if not logged in", func(t *testing.T) {
		initAuthority(ctx, t, tinit.DB(t))
		targetUser := signUpTestUser(t)
		res, err := http.DefaultClient.Post(rootUrl+"/auth/admin/authority", "application/json", strings.NewReader(fmt.Sprintf(sampleJsonBody, targetUser.UserID)))
		require.NoError(t, err)

		require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("return 403 if not authorized", func(t *testing.T) {
		initAuthority(ctx, t, tinit.DB(t))

		tokens, err := jwtResolver.CreateToken("notImportant", []string{"AUTHORITY_USER"}, time.Now()) //not AUTHORITY_ADMIN
		require.NoError(t, err)

		targetUser := signUpTestUser(t)
		req, err := http.NewRequest("POST", rootUrl+"/auth/admin/authority", strings.NewReader(fmt.Sprintf(sampleJsonBody, targetUser.UserID)))
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusForbidden, res.StatusCode)
	})

	t.Run("return 201 if successfully added", func(t *testing.T) {
		initAuthority(ctx, t, tinit.DB(t))

		tokens, err := jwtResolver.CreateToken("notImportant", []string{domain.AuthorityAdmin}, time.Now())
		require.NoError(t, err)

		targetUser := signUpTestUser(t)
		req, err := http.NewRequest("POST", rootUrl+"/auth/admin/authority", strings.NewReader(fmt.Sprintf(sampleJsonBody, targetUser.UserID)))
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, http.StatusCreated, res.StatusCode)

		userAuthorities, err := userService.FindUserAuthorities(ctx, targetUser.UserID)
		require.NoError(t, err)
		require.Len(t, userAuthorities, 2)
		require.Equal(t, "AUTHORITY_USER", userAuthorities[0].AuthorityName)
		require.WithinDuration(t, time.Now().Add(720*time.Hour), *userAuthorities[0].ExpiryDate, 1*time.Second)
		require.Equal(t, "AUTHORITY_GUEST", userAuthorities[1].AuthorityName)
		require.Nil(t, userAuthorities[1].ExpiryDate)
	})
}

func initRootUrl(t *testing.T) string {
	envVars := tinit.InitEnvVars(t)

	return "http://localhost:" + strconv.Itoa(envVars.ApiPort)
}

func initJwtResolver(t *testing.T) *jwtresolver.JwtResolver {
	envVars := tinit.InitEnvVars(t)
	return jwtresolver.NewJwtResolver(envVars.SecretKey)
}
