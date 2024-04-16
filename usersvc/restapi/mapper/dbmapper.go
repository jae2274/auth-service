package mapper

import (
	"context"
	"database/sql"
	"userService/models"
	"userService/usersvc/common/domain"

	"github.com/jae2274/goutils/terr"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func SaveUser(ctx context.Context, exec boil.ContextExecutor, authorizedBy domain.AuthorizedBy, authorizedID, email string) error {
	user := models.User{
		AuthorizedBy: string(authorizedBy),
		AuthorizedID: authorizedID,
		Email:        email,
	}

	err := user.Insert(ctx, exec, boil.Infer())
	if err != nil {
		return terr.Wrap(err)
	}

	return nil
}

func FindByAuthorized(ctx context.Context, exec boil.ContextExecutor, authorizedType domain.AuthorizedBy, authorizedID string) (*models.User, error) {

	user, err := models.Users(qm.Where(models.UserColumns.AuthorizedBy+"=?", authorizedType), qm.And(models.UserColumns.AuthorizedID+"=?", authorizedID)).One(ctx, exec)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, terr.Wrap(err)
	}

	return user, nil
}
