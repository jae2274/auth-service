package mapper

import (
	"context"
	"database/sql"
	"userService/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func GetUserEMails(ctx context.Context, exec boil.ContextExecutor, userIds []int) ([]*models.User, error) {
	if len(userIds) == 0 {
		return []*models.User{}, nil
	}

	users, err := models.Users(qm.Where(models.UserColumns.UserID+" IN ?", userIds, qm.And(models.UserColumns.AgreeMail+"=?", true))).All(ctx, exec)

	if err != nil {
		if err == sql.ErrNoRows {
			return []*models.User{}, nil
		}
	}

	return users, nil
}
