package service

import (
	"database/sql"
	"fmt"
	"userService/usersvc/common/domain"
	"userService/usersvc/common/entity"
	"userService/usersvc/restapi/mapper"
)

type UserService interface {
	GetUser(authorizedBy domain.AuthorizedBy, authorizedID string) (*domain.User, error)
	SaveUser(authorizedBy domain.AuthorizedBy, authorizedID, email string) error
}

type UserServiceImpl struct {
	mysqlDB *sql.DB
}

func NewUserService(mysqlDB *sql.DB) UserService {
	return &UserServiceImpl{
		mysqlDB: mysqlDB,
	}
}

func (u *UserServiceImpl) GetUser(authorizedBy domain.AuthorizedBy, authorizedID string) (*domain.User, error) {
	tx, err := u.mysqlDB.Begin()
	if err != nil {
		return nil, err
	}

	user, err := u.getUser(tx, authorizedBy, authorizedID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserServiceImpl) getUser(tx *sql.Tx, authorizedBy domain.AuthorizedBy, authorizedID string) (*domain.User, error) {
	user, err := mapper.FindByAuthorized(tx, authorizedBy, authorizedID)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, nil
	}

	userRoles, err := mapper.FindAllUserRoles(tx, user.UserID)
	if err != nil {
		return nil, err
	}

	roles := make([]string, len(userRoles))
	for i, role := range userRoles {
		roles[i] = role.RoleName
	}

	return &domain.User{
		UserID:       fmt.Sprintf("%d", user.UserID),
		AuthorizedBy: user.AuthorizedBy,
		AuthorizedID: user.AuthorizedID,
		Email:        user.Email,
		CreateDate:   user.CreateDate,
		Roles:        roles,
	}, nil
}

func (u *UserServiceImpl) SaveUser(authorizedBy domain.AuthorizedBy, authorizedID, email string) error {
	tx, err := u.mysqlDB.Begin()
	if err != nil {
		return err
	}

	if err := u.saveUser(tx, authorizedBy, authorizedID, email); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (u *UserServiceImpl) saveUser(tx *sql.Tx, authorizedBy domain.AuthorizedBy, authorizedID, email string) error {
	userVO := entity.UserVO{
		AuthorizedBy: authorizedBy,
		AuthorizedID: authorizedID,
		Email:        email,
	}

	return mapper.SaveUser(tx, userVO)
}
