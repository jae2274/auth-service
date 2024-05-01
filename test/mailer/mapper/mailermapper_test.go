package mapper

import (
	"context"
	"testing"
	"userService/test/tinit"
	"userService/test/tutils"
	"userService/usersvc/common/domain"
	mailMapper "userService/usersvc/mailer/mapper"
	"userService/usersvc/models"
	restapiMapper "userService/usersvc/restapi/mapper"

	"github.com/stretchr/testify/require"
)

func TestMailerDBMapper(t *testing.T) {

	t.Run("return empty when nothing saved", func(t *testing.T) {
		sqlDB := tinit.DB(t)
		defer sqlDB.Close()

		ctx := context.Background()
		userEMails, err := mailMapper.GetUserEMails(ctx, sqlDB, []int{})
		require.NoError(t, err)
		require.Len(t, userEMails, 0)
	})

	t.Run("return empty user if not existed", func(t *testing.T) {
		sqlDB := tinit.DB(t)
		defer sqlDB.Close()

		ctx := context.Background()
		userEMails, err := mailMapper.GetUserEMails(ctx, sqlDB, []int{1})
		require.NoError(t, err)
		require.Len(t, userEMails, 0)
	})

	t.Run("return user if existed", func(t *testing.T) {
		sqlDB := tinit.DB(t)
		defer sqlDB.Close()

		ctx := context.Background()
		willSavedUser := tutils.NewUser(1)
		saved, err := restapiMapper.SaveUser(ctx, sqlDB, domain.AuthorizedBy(willSavedUser.AuthorizedBy), willSavedUser.AuthorizedID, willSavedUser.Email, willSavedUser.Name)
		require.NoError(t, err)

		userEMails, err := mailMapper.GetUserEMails(ctx, sqlDB, []int{saved.UserID})
		require.NoError(t, err)
		require.Len(t, userEMails, 1)
		require.Equal(t, willSavedUser.AuthorizedBy, userEMails[0].AuthorizedBy)
		require.Equal(t, willSavedUser.AuthorizedID, userEMails[0].AuthorizedID)
		require.Equal(t, willSavedUser.Email, userEMails[0].Email)
		require.Equal(t, willSavedUser.Name, userEMails[0].Name)
	})

	t.Run("return users if all existed", func(t *testing.T) {
		sqlDB := tinit.DB(t)
		defer sqlDB.Close()

		ctx := context.Background()
		willSavedUsers := []*models.User{tutils.NewUser(1), tutils.NewUser(2), tutils.NewUser(3), tutils.NewUser(4)}

		userIds := make([]int, len(willSavedUsers))

		for i, willSavedUser := range willSavedUsers {
			saved, err := restapiMapper.SaveUser(ctx, sqlDB, domain.AuthorizedBy(willSavedUser.AuthorizedBy), willSavedUser.AuthorizedID, willSavedUser.Email, willSavedUser.Name)
			require.NoError(t, err)
			userIds[i] = saved.UserID
		}

		findedUsers, err := mailMapper.GetUserEMails(ctx, sqlDB, userIds)
		require.NoError(t, err)
		require.Len(t, findedUsers, len(willSavedUsers))
		for i, u := range findedUsers {
			require.Equal(t, willSavedUsers[i].AuthorizedBy, u.AuthorizedBy)
			require.Equal(t, willSavedUsers[i].AuthorizedID, u.AuthorizedID)
			require.Equal(t, willSavedUsers[i].Email, u.Email)
			require.Equal(t, willSavedUsers[i].Name, u.Name)
		}
	})

	t.Run("return some existed users", func(t *testing.T) {
		sqlDB := tinit.DB(t)
		defer sqlDB.Close()

		ctx := context.Background()
		willSavedUsers := []*models.User{tutils.NewUser(1), tutils.NewUser(2), tutils.NewUser(3), tutils.NewUser(4)}

		userIds := make([]int, len(willSavedUsers))

		for i, willSavedUser := range willSavedUsers {
			saved, err := restapiMapper.SaveUser(ctx, sqlDB, domain.AuthorizedBy(willSavedUser.AuthorizedBy), willSavedUser.AuthorizedID, willSavedUser.Email, willSavedUser.Name)
			require.NoError(t, err)
			userIds[i] = saved.UserID
		}

		findedUsers, err := mailMapper.GetUserEMails(ctx, sqlDB, []int{userIds[0], userIds[2], 9999})
		require.NoError(t, err)
		require.Len(t, findedUsers, 2)
		for i, u := range findedUsers {
			require.Equal(t, willSavedUsers[i*2].AuthorizedBy, u.AuthorizedBy)
			require.Equal(t, willSavedUsers[i*2].AuthorizedID, u.AuthorizedID)
			require.Equal(t, willSavedUsers[i*2].Email, u.Email)
			require.Equal(t, willSavedUsers[i*2].Name, u.Name)
		}
	})
}
