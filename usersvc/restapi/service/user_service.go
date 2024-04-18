package service

import (
	"context"
	"database/sql"
	"fmt"
	"userService/usersvc/common/domain"
	"userService/usersvc/restapi/mapper"
)

type UserService interface {
	GetUser(ctx context.Context, authorizedBy domain.AuthorizedBy, authorizedID string) (*domain.User, bool, error)
	SaveUser(ctx context.Context, authorizedBy domain.AuthorizedBy, authorizedID, email string) error
}

type UserServiceImpl struct {
	mysqlDB *sql.DB
}

func NewUserService(mysqlDB *sql.DB) UserService {
	return &UserServiceImpl{
		mysqlDB: mysqlDB,
	}
}

func (u *UserServiceImpl) GetUser(ctx context.Context, authorizedBy domain.AuthorizedBy, authorizedID string) (*domain.User, bool, error) {
	tx, err := u.mysqlDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, false, err
	}

	user, isExisted, err := u.getUser(ctx, tx, authorizedBy, authorizedID)
	if err != nil {
		tx.Rollback()
		return nil, false, err
	}

	if err := tx.Commit(); err != nil {
		return nil, false, err
	}

	return user, isExisted, nil
}

func (u *UserServiceImpl) getUser(ctx context.Context, tx *sql.Tx, authorizedBy domain.AuthorizedBy, authorizedID string) (*domain.User, bool, error) {
	user, isExisted, err := mapper.FindUserByAuthorized(ctx, tx, authorizedBy, authorizedID)
	if err != nil {
		return nil, false, err
	} else if !isExisted {
		return nil, false, nil
	}

	userRoles, err := user.UserRoles().All(ctx, tx)
	if err != nil {
		return nil, false, err
	}

	roles := make([]string, len(userRoles))
	for i, role := range userRoles {
		roles[i] = role.RoleName
	}

	return &domain.User{
		UserID:       fmt.Sprintf("%d", user.UserID),
		AuthorizedBy: domain.AuthorizedBy(user.AuthorizedBy),
		AuthorizedID: user.AuthorizedID,
		Email:        user.Email,
		CreateDate:   user.CreateDate,
		Roles:        roles,
	}, true, nil
}

func (u *UserServiceImpl) SaveUser(ctx context.Context, authorizedBy domain.AuthorizedBy, authorizedID, email string) error {
	return mapper.SaveUser(ctx, u.mysqlDB, authorizedBy, authorizedID, email)
}
