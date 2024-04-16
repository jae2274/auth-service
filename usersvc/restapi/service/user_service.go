package service

import (
	"context"
	"database/sql"
	"fmt"
	"userService/usersvc/common/domain"
	"userService/usersvc/restapi/mapper"
)

type UserService interface {
	GetUser(ctx context.Context, authorizedBy domain.AuthorizedBy, authorizedID string) (*domain.User, error)
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

func (u *UserServiceImpl) GetUser(ctx context.Context, authorizedBy domain.AuthorizedBy, authorizedID string) (*domain.User, error) {
	tx, err := u.mysqlDB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	user, err := u.getUser(ctx, tx, authorizedBy, authorizedID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserServiceImpl) getUser(ctx context.Context, tx *sql.Tx, authorizedBy domain.AuthorizedBy, authorizedID string) (*domain.User, error) {
	user, err := mapper.FindByAuthorized(ctx, tx, authorizedBy, authorizedID)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, nil
	}

	userRoles, err := user.UserRoles().All(ctx, tx)
	if err != nil {
		return nil, err
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
	}, nil
}

func (u *UserServiceImpl) SaveUser(ctx context.Context, authorizedBy domain.AuthorizedBy, authorizedID, email string) error {
	return mapper.SaveUser(ctx, u.mysqlDB, authorizedBy, authorizedID, email)
}
