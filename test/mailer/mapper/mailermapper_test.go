package mapper

import (
	"context"
	"database/sql"
	"testing"
	"userService/models"
	"userService/test/tinit"
	"userService/test/tutils"
	"userService/usersvc/common/domain"
	mailMapper "userService/usersvc/mailer/mapper"
	restapiMapper "userService/usersvc/restapi/mapper"

	"github.com/stretchr/testify/require"
)

func TestMailerDBMapper(t *testing.T) {
	willSavedUserVO := models.User{
		AuthorizedBy: string(domain.GOOGLE),
		AuthorizedID: "test",
		Email:        "test@mail.com",
	}

	t.Run("No one agree to mail", func(t *testing.T) {
		sqlDB := tinit.DB(t)
		defer sqlDB.Close()

		ctx := context.Background()
		tutils.TxCommit(t, sqlDB, func(tx *sql.Tx) {
			err := restapiMapper.SaveUser(ctx, tx, domain.AuthorizedBy(willSavedUserVO.AuthorizedBy), willSavedUserVO.AuthorizedID, willSavedUserVO.Email)
			require.NoError(t, err)

			userEMails, err := mailMapper.GetUserEMails(ctx, tx, []int{1})
			require.NoError(t, err)
			require.Len(t, userEMails, 0)
		})
	})
}
