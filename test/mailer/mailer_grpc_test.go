package mailer

import (
	"context"
	"strconv"
	"testing"
	"userService/test/tinit"
	"userService/test/tutils"
	"userService/usersvc/common/domain"
	"userService/usersvc/mailer/mailer_grpc"
	"userService/usersvc/models"
	"userService/usersvc/restapi/mapper"

	"github.com/stretchr/testify/require"
)

func TestMailerGrpc(t *testing.T) {

	cancelFunc := tinit.RunTestApp(t)
	defer cancelFunc()
	t.Run("return empty", func(t *testing.T) {
		tinit.DB(t)
		client := initMailerClient(t)

		ctx := context.Background()
		users, err := client.GetUserEmails(ctx, &mailer_grpc.GetUserEmailsRequest{
			UserIds: []string{"1"},
		})
		require.NoError(t, err)
		require.Len(t, users.UserEmails, 0)
	})

	t.Run("return emails", func(t *testing.T) {
		ctx := context.Background()
		//given
		sqlDB := tinit.DB(t)
		willSavedUsers := []*models.User{tutils.NewUser(1), tutils.NewUser(2), tutils.NewUser(3)}
		savedUsers := make([]*models.User, len(willSavedUsers))

		for i, willSavedUser := range willSavedUsers {
			saved, err := mapper.SaveUser(ctx, sqlDB, domain.AuthorizedBy(willSavedUser.AuthorizedBy), willSavedUser.AuthorizedID, willSavedUser.Email, willSavedUser.Name)
			require.NoError(t, err)
			savedUsers[i] = saved
		}

		client := initMailerClient(t)

		//when
		users, err := client.GetUserEmails(ctx, &mailer_grpc.GetUserEmailsRequest{
			UserIds: []string{"1", "2", "3", "4"},
		})
		require.NoError(t, err)

		//then
		require.Len(t, users.UserEmails, 3)

		for i, user := range users.UserEmails {
			intUserId, err := strconv.Atoi(user.UserId)
			require.NoError(t, err)
			require.Equal(t, savedUsers[i].UserID, intUserId)
			require.Equal(t, savedUsers[i].Email, user.Email)
		}
	})
}

func initMailerClient(t *testing.T) mailer_grpc.UserClient {
	envVars := tinit.InitEnvVars(t)
	conn := tinit.InitGrpcClient(t, envVars.MailerGrpcPort)

	return mailer_grpc.NewUserClient(conn)
}
