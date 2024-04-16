package mapper

import (
	"database/sql"
	"testing"
	"userService/test/tinit"
	"userService/test/tutils"
	"userService/usersvc/common/domain"
	"userService/usersvc/common/entity"
	mailMapper "userService/usersvc/mailer/mapper"
	restapiMapper "userService/usersvc/restapi/mapper"

	"github.com/stretchr/testify/require"
)

func TestMailerDBMapper(t *testing.T) {
	willSavedUserVO := entity.UserVO{
		AuthorizedBy: domain.GOOGLE,
		AuthorizedID: "test",
		Email:        "test@mail.com",
	}

	t.Run("No one agree to mail", func(t *testing.T) {
		sqlDB := tinit.DB(t)
		tutils.TxCommit(t, sqlDB, func(tx *sql.Tx) {
			err := restapiMapper.SaveUser(tx, willSavedUserVO)
			require.NoError(t, err)

			userEMails, err := mailMapper.GetUserEMails(tx, []int64{1})
			require.NoError(t, err)
			require.Len(t, userEMails, 0)
		})
	})
}
