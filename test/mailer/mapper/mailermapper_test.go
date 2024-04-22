package mapper

import (
	"context"
	"database/sql"
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
	willSavedUserVO := models.User{
		AuthorizedBy: string(domain.GOOGLE),
		AuthorizedID: "test",
		Name:         "testName",
		Email:        "test@mail.com",
	}

	t.Run("No one agree to mail", func(t *testing.T) {
		sqlDB := tinit.DB(t)
		defer sqlDB.Close()

		ctx := context.Background()
		tutils.TxCommit(t, sqlDB, func(tx *sql.Tx) {
			u, err := restapiMapper.SaveUser(ctx, tx, domain.AuthorizedBy(willSavedUserVO.AuthorizedBy), willSavedUserVO.AuthorizedID, willSavedUserVO.Email, willSavedUserVO.Name)
			require.NoError(t, err)
			require.NotZero(t, u.UserID)

			userEMails, err := mailMapper.GetUserEMails(ctx, tx, []int{1})
			require.NoError(t, err)
			require.Len(t, userEMails, 0)
		})
	})
}
