package controller

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/jae2274/auth-service/auth_service/common/domain"
	"github.com/jae2274/auth-service/auth_service/common/mysqldb"
	"github.com/jae2274/auth-service/auth_service/models"
	"github.com/jae2274/auth-service/auth_service/restapi/ctrlr/dto"
	"github.com/jae2274/auth-service/auth_service/restapi/jwtresolver"
	"github.com/jae2274/auth-service/auth_service/restapi/ooauth"
	"github.com/jae2274/auth-service/auth_service/restapi/service"
	"github.com/jae2274/auth-service/test/tinit"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

func newAuthorities() []*models.Authority {
	return []*models.Authority{
		{AuthorityCode: domain.AuthorityAdmin, AuthorityName: "관리자", Summary: "관리자 권한"},
		{AuthorityCode: "AUTHORITY_USER", AuthorityName: "사용자", Summary: "사용자 권한"},
		{AuthorityCode: "AUTHORITY_GUEST", AuthorityName: "손님", Summary: "게스트 권한"},
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

func initAdminUser(ctx context.Context, t *testing.T, db *sql.DB) *jwtresolver.CustomClaims {
	user, err := mysqldb.WithTransaction(ctx, db, func(tx *sql.Tx) (*models.User, error) {
		user, err := service.SignUp(ctx, tx, &ooauth.UserInfo{
			AuthorizedID: "authorizedID",
			AuthorizedBy: "GOOGLE",
			Email:        "testEmail@testmail.net",
			Username:     "testUsername",
		}, []*dto.UserAgreementReq{})

		return user, err
	})

	require.NoError(t, err)
	return &jwtresolver.CustomClaims{
		UserId:       strconv.Itoa(user.UserID),
		AuthorizedBy: domain.AuthorizedBy(user.AuthorizedBy),
		AuthorizedID: user.AuthorizedID,
	}
}

// admin api는 jwt의 토큰을 기반으로 동작한다.
// jwt에 admin 권한이 존재하는지, 정지되거나 탈퇴한 계정인지 확인을 위해 DB조회를 할 뿐, DB의 권한 정보는 api호출에 영향을 미치지 않는다.
func TestAdminController(t *testing.T) {
	ctx := context.Background()

	cancelFunc := tinit.RunTestApp(t)
	defer cancelFunc()

	rootUrl := initRootUrl(t)
	jwtResolver := initJwtResolver(t)

	signUpTestUser := func(t *testing.T, db *sql.DB) *models.User {
		user, err := mysqldb.WithTransaction(ctx, db, func(tx *sql.Tx) (*models.User, error) {
			return service.SignUp(ctx, tx, &ooauth.UserInfo{AuthorizedBy: domain.GOOGLE, AuthorizedID: "authId", Email: "targetUser@test.com", Username: "target"}, nil)
		})
		require.NoError(t, err)

		return user
	}

	addSampleJsonBody := `
	{
		"userId": %d,
		"authorities": [
		  {
			"authorityCode": "AUTHORITY_USER",
			"expiryDurationMS": 2592000000
		  },
		  {
			"authorityCode": "AUTHORITY_GUEST"
		  }
		]
	  }
	`

	removeSampleJsonBody := `
	{
		"userId": %d,
		"authorityCode": "AUTHORITY_USER"
	  }
	`

	adminAuthority := "/auth/admin/authority"
	adminTicket := "/auth/admin/ticket"
	t.Run("GetAllAuthorities", func(t *testing.T) {

		t.Run("return 401 if not logged in", func(t *testing.T) {
			db := tinit.DB(t)
			initAuthority(ctx, t, db)
			res, err := http.DefaultClient.Get(rootUrl + adminAuthority)
			require.NoError(t, err)

			require.Equal(t, http.StatusUnauthorized, res.StatusCode)
		})

		t.Run("return 403 if not authorized", func(t *testing.T) {
			db := tinit.DB(t)
			initAuthority(ctx, t, db)

			admin := initAdminUser(ctx, t, db)                                                                                                   //admin api 호출시
			tokens, err := jwtResolver.CreateToken(admin.UserId, admin.AuthorizedBy, admin.AuthorizedID, []string{"AUTHORITY_USER"}, time.Now()) //not AUTHORITY_ADMIN
			require.NoError(t, err)

			req, err := http.NewRequest("GET", rootUrl+adminAuthority, nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			require.Equal(t, http.StatusForbidden, res.StatusCode)
		})

		t.Run("return 200 with authorities", func(t *testing.T) {
			db := tinit.DB(t)
			authorities := initAuthority(ctx, t, db)

			admin := initAdminUser(ctx, t, db)
			tokens, err := jwtResolver.CreateToken(admin.UserId, admin.AuthorizedBy, admin.AuthorizedID, []string{domain.AuthorityAdmin}, time.Now())
			require.NoError(t, err)

			req, err := http.NewRequest("GET", rootUrl+adminAuthority, nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			require.Equal(t, http.StatusOK, res.StatusCode)

			var resObj dto.GetAllAuthoritiesResponse
			json.NewDecoder(res.Body).Decode(&resObj)

			require.Len(t, resObj.Authorities, len(authorities)-1) //관리자 권한은 제외
			for i, authority := range authorities[1:] {            //관리자 권한은 제외
				require.Equal(t, authority.AuthorityCode, resObj.Authorities[i].AuthorityCode)
				require.Equal(t, authority.AuthorityName, resObj.Authorities[i].AuthorityName)
				require.Equal(t, authority.Summary, resObj.Authorities[i].Summary)
			}
		})
	})

	t.Run("AddAuthority", func(t *testing.T) {

		t.Run("return 401 if not logged in", func(t *testing.T) {
			db := tinit.DB(t)
			initAuthority(ctx, t, db)
			targetUser := signUpTestUser(t, db)
			res, err := http.DefaultClient.Post(rootUrl+adminAuthority, "application/json", strings.NewReader(fmt.Sprintf(addSampleJsonBody, targetUser.UserID)))
			require.NoError(t, err)

			require.Equal(t, http.StatusUnauthorized, res.StatusCode)
		})

		t.Run("return 403 if not authorized", func(t *testing.T) {
			db := tinit.DB(t)
			initAuthority(ctx, t, db)

			admin := initAdminUser(ctx, t, db)
			tokens, err := jwtResolver.CreateToken(admin.UserId, admin.AuthorizedBy, admin.AuthorizedID, []string{"AUTHORITY_USER"}, time.Now()) //not AUTHORITY_ADMIN
			require.NoError(t, err)

			targetUser := signUpTestUser(t, db)
			req, err := http.NewRequest("POST", rootUrl+adminAuthority, strings.NewReader(fmt.Sprintf(addSampleJsonBody, targetUser.UserID)))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			require.Equal(t, http.StatusForbidden, res.StatusCode)
		})

		t.Run("return 201 if successfully added", func(t *testing.T) {
			db := tinit.DB(t)
			initAuthority(ctx, t, db)

			admin := initAdminUser(ctx, t, db)
			tokens, err := jwtResolver.CreateToken(admin.UserId, admin.AuthorizedBy, admin.AuthorizedID, []string{domain.AuthorityAdmin}, time.Now())
			require.NoError(t, err)

			targetUser := signUpTestUser(t, db)
			req, err := http.NewRequest("POST", rootUrl+adminAuthority, strings.NewReader(fmt.Sprintf(addSampleJsonBody, targetUser.UserID)))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			require.Equal(t, http.StatusCreated, res.StatusCode)

			userAuthorities, err := service.FindValidUserAuthorities(ctx, db, targetUser.UserID)
			require.NoError(t, err)
			require.Len(t, userAuthorities, 2)
			require.Equal(t, "AUTHORITY_USER", userAuthorities[0].AuthorityCode)
			require.WithinDuration(t, time.Now().Add(720*time.Hour).UTC(), time.UnixMilli(*(userAuthorities[0].ExpiryUnixMilli)).UTC(), 1*time.Second)
			require.Equal(t, "AUTHORITY_GUEST", userAuthorities[1].AuthorityCode)
			require.Nil(t, userAuthorities[1].ExpiryUnixMilli)
		})
	})

	t.Run("RemoveAuthority", func(t *testing.T) {
		t.Run("return 401 if not logged in", func(t *testing.T) {
			db := tinit.DB(t)
			initAuthority(ctx, t, db)
			targetUser := signUpTestUser(t, db)
			req, err := http.NewRequest("DELETE", rootUrl+adminAuthority, strings.NewReader(fmt.Sprintf(removeSampleJsonBody, targetUser.UserID)))
			require.NoError(t, err)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			require.Equal(t, http.StatusUnauthorized, res.StatusCode)
		})

		t.Run("return 403 if not authorized", func(t *testing.T) {
			db := tinit.DB(t)
			initAuthority(ctx, t, db)

			admin := initAdminUser(ctx, t, db)
			tokens, err := jwtResolver.CreateToken(admin.UserId, admin.AuthorizedBy, admin.AuthorizedID, []string{"AUTHORITY_USER"}, time.Now()) //not AUTHORITY_ADMIN
			require.NoError(t, err)

			targetUser := signUpTestUser(t, db)
			req, err := http.NewRequest("DELETE", rootUrl+adminAuthority, strings.NewReader(fmt.Sprintf(removeSampleJsonBody, targetUser.UserID)))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			require.Equal(t, http.StatusForbidden, res.StatusCode)
		})

		t.Run("return 204 if successfully removed", func(t *testing.T) {
			db := tinit.DB(t)
			initAuthority(ctx, t, db)

			admin := initAdminUser(ctx, t, db)
			tokens, err := jwtResolver.CreateToken(admin.UserId, admin.AuthorizedBy, admin.AuthorizedID, []string{domain.AuthorityAdmin}, time.Now())
			require.NoError(t, err)

			targetUser := signUpTestUser(t, db)

			req, err := http.NewRequest("POST", rootUrl+adminAuthority, strings.NewReader(fmt.Sprintf(addSampleJsonBody, targetUser.UserID)))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)
			_, err = http.DefaultClient.Do(req)
			require.NoError(t, err)

			userAuthorities, err := service.FindValidUserAuthorities(ctx, db, targetUser.UserID)
			require.NoError(t, err)
			require.Len(t, userAuthorities, 2)

			req, err = http.NewRequest("DELETE", rootUrl+adminAuthority, strings.NewReader(fmt.Sprintf(removeSampleJsonBody, targetUser.UserID)))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			require.Equal(t, http.StatusNoContent, res.StatusCode)

			userAuthorities, err = service.FindValidUserAuthorities(ctx, db, targetUser.UserID)
			require.NoError(t, err)
			require.Len(t, userAuthorities, 1)
		})
	})

	createTicketReq := `
	{
		"useableCount": 1,
		"ticketName": "testTicket",
		"ticketAuthorities": [
			{
				"authorityCode": "AUTHORITY_USER",
				"expiryDurationMS": 2592000000
			},
			{
				"authorityCode": "AUTHORITY_GUEST"
			}
		]
	}
	`
	t.Run("CreateTicket", func(t *testing.T) {

		t.Run("return 401 if not logged in", func(t *testing.T) {
			db := tinit.DB(t)
			initAuthority(ctx, t, db)
			res, err := http.DefaultClient.Post(rootUrl+adminTicket, "application/json", strings.NewReader(createTicketReq))
			require.NoError(t, err)

			require.Equal(t, http.StatusUnauthorized, res.StatusCode)
		})

		t.Run("return 403 if not authorized", func(t *testing.T) {
			db := tinit.DB(t)
			initAuthority(ctx, t, db)

			admin := initAdminUser(ctx, t, db)
			tokens, err := jwtResolver.CreateToken(admin.UserId, admin.AuthorizedBy, admin.AuthorizedID, []string{"AUTHORITY_USER"}, time.Now()) //not AUTHORITY_ADMIN
			require.NoError(t, err)

			req, err := http.NewRequest("POST", rootUrl+adminTicket, strings.NewReader(createTicketReq))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			require.Equal(t, http.StatusForbidden, res.StatusCode)
		})

		t.Run("return 201 if successfully created", func(t *testing.T) {
			db := tinit.DB(t)
			initAuthority(ctx, t, db)

			admin := initAdminUser(ctx, t, db)
			tokens, err := jwtResolver.CreateToken(admin.UserId, admin.AuthorizedBy, admin.AuthorizedID, []string{domain.AuthorityAdmin}, time.Now())
			require.NoError(t, err)

			req, err := http.NewRequest("POST", rootUrl+adminTicket, strings.NewReader(createTicketReq))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			require.Equal(t, http.StatusCreated, res.StatusCode)

			var resObj *dto.Ticket
			json.NewDecoder(res.Body).Decode(&resObj)

			require.Len(t, resObj.TicketAuthorities, 2)
			require.NotEmpty(t, resObj.TicketId)
			require.Equal(t, "AUTHORITY_USER", resObj.TicketAuthorities[0].AuthorityCode)
			require.Equal(t, 720*time.Hour/time.Millisecond, time.Duration(*resObj.TicketAuthorities[0].ExpiryDurationMS))
			require.Equal(t, "AUTHORITY_GUEST", resObj.TicketAuthorities[1].AuthorityCode)
			require.Nil(t, resObj.TicketAuthorities[1].ExpiryDurationMS)

			_, isExisted, err := service.GetTicketInfo(ctx, db, resObj.TicketId)
			require.NoError(t, err)
			require.True(t, isExisted)
		})
	})

	t.Run("GetAllTickets", func(t *testing.T) {
		t.Run("return 401 if not logged in", func(t *testing.T) {
			db := tinit.DB(t)
			initAuthority(ctx, t, db)
			res, err := http.DefaultClient.Get(rootUrl + adminTicket)
			require.NoError(t, err)

			require.Equal(t, http.StatusUnauthorized, res.StatusCode)
		})

		t.Run("return 403 if not authorized", func(t *testing.T) {
			db := tinit.DB(t)
			initAuthority(ctx, t, db)

			admin := initAdminUser(ctx, t, db)
			tokens, err := jwtResolver.CreateToken(admin.UserId, admin.AuthorizedBy, admin.AuthorizedID, []string{"AUTHORITY_USER"}, time.Now()) //not AUTHORITY_ADMIN
			require.NoError(t, err)

			req, err := http.NewRequest("GET", rootUrl+adminTicket, nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			require.Equal(t, http.StatusForbidden, res.StatusCode)
		})

		t.Run("return 200 with tickets", func(t *testing.T) {
			db := tinit.DB(t)
			initAuthority(ctx, t, db)

			admin := initAdminUser(ctx, t, db)
			tokens, err := jwtResolver.CreateToken(admin.UserId, admin.AuthorizedBy, admin.AuthorizedID, []string{domain.AuthorityAdmin}, time.Now())
			require.NoError(t, err)

			req, err := http.NewRequest("POST", rootUrl+adminTicket, strings.NewReader(createTicketReq))
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

			_, err = http.DefaultClient.Do(req)
			require.NoError(t, err)

			req, err = http.NewRequest("GET", rootUrl+adminTicket, nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+tokens.AccessToken)

			res, err := http.DefaultClient.Do(req)
			require.NoError(t, err)

			require.Equal(t, http.StatusOK, res.StatusCode)
			var resObj dto.GetAllTicketsResponse
			json.NewDecoder(res.Body).Decode(&resObj)

			require.Len(t, resObj.Tickets, 1)
			ticketRes := resObj.Tickets[0]
			require.Len(t, ticketRes.TicketAuthorities, 2)
			require.NotEmpty(t, ticketRes.TicketId)
			require.Equal(t, "AUTHORITY_USER", ticketRes.TicketAuthorities[0].AuthorityCode)
			require.Equal(t, 720*time.Hour/time.Millisecond, time.Duration(*ticketRes.TicketAuthorities[0].ExpiryDurationMS))
			require.Equal(t, "AUTHORITY_GUEST", ticketRes.TicketAuthorities[1].AuthorityCode)
			require.Nil(t, ticketRes.TicketAuthorities[1].ExpiryDurationMS)
		})
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
