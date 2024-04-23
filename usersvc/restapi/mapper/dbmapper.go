package mapper

import (
	"context"
	"database/sql"
	"userService/usersvc/common/domain"
	"userService/usersvc/models"

	"github.com/jae2274/goutils/terr"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func SaveUser(ctx context.Context, exec boil.ContextExecutor, authorizedBy domain.AuthorizedBy, authorizedID, email, username string) (*models.User, error) {
	user := &models.User{
		AuthorizedBy: string(authorizedBy),
		AuthorizedID: authorizedID,
		Email:        email,
		Name:         username,
	}

	err := user.Insert(ctx, exec, boil.Infer())
	if err != nil {
		return nil, terr.Wrap(err)
	}

	return user, nil
}

func FindUserByAuthorized(ctx context.Context, exec boil.ContextExecutor, authorizedType domain.AuthorizedBy, authorizedID string) (*models.User, bool, error) {

	user, err := models.Users(qm.Where(models.UserColumns.AuthorizedBy+"=?", authorizedType), qm.And(models.UserColumns.AuthorizedID+"=?", authorizedID), qm.Load(models.UserRels.UserRoles)).One(ctx, exec)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, terr.Wrap(err)
	}

	return user, true, nil
}

func FindUserAgreements(ctx context.Context, exec boil.ContextExecutor, userID int, isAgree bool) ([]*models.UserAgreement, error) {
	userAgreements, err := models.UserAgreements(qm.Where(models.UserAgreementColumns.UserID+"=?", userID), qm.And(models.UserAgreementColumns.IsAgree+"=?", isAgree)).All(ctx, exec)
	if err != nil {
		return nil, terr.Wrap(err)
	}

	return userAgreements, nil
}

func FindAllAgreement(ctx context.Context, exec boil.ContextExecutor) ([]*models.Agreement, error) {
	agreements, err := models.Agreements(qm.OrderBy(models.AgreementColumns.Priority)).All(ctx, exec)
	if err != nil {
		return nil, terr.Wrap(err)
	}

	return agreements, nil
}

func FindRequiredAgreements(ctx context.Context, exec boil.ContextExecutor) ([]*models.Agreement, error) {
	agreements, err := models.Agreements(qm.Where(models.AgreementColumns.IsRequired+"=?", true), qm.OrderBy(models.AgreementColumns.Priority)).All(ctx, exec)
	if err != nil {
		return nil, terr.Wrap(err)
	}

	return agreements, nil
}
