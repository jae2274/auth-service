package mapper

import (
	"context"
	"database/sql"

	"github.com/jae2274/auth-service/auth_service/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func GetUserEMails(ctx context.Context, exec boil.ContextExecutor, userIds []int) ([]*models.User, error) {
	if len(userIds) == 0 {
		return []*models.User{}, nil
	}
	convertedUserIds := make([]interface{}, len(userIds))
	for i, v := range userIds {
		convertedUserIds[i] = v
	}

	users, err := models.Users(qm.WhereIn(models.UserColumns.UserID+" IN ?", convertedUserIds...)).All(ctx, exec)

	if err != nil {
		if err == sql.ErrNoRows {
			return users, nil
		}
		return users, err
	}

	return users, nil
}
