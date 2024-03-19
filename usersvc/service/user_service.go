package service

import (
	"fmt"
	"userService/usersvc/domain"
	"userService/usersvc/entity"
	"userService/usersvc/mysqldb"
)

type UserService interface {
	GetUser(authorizedBy domain.AuthorizedBy, authorizedID string) (*domain.User, error)
	SaveUser(authorizedBy domain.AuthorizedBy, authorizedID, email string) error
}

type UserServiceImpl struct {
	UserRepo mysqldb.DBMapper
}

func NewUserService(userRepo mysqldb.DBMapper) UserService {
	return &UserServiceImpl{
		UserRepo: userRepo,
	}
}

func (u *UserServiceImpl) GetUser(authorizedBy domain.AuthorizedBy, authorizedID string) (*domain.User, error) {
	user, err := u.UserRepo.FindByAuthorized(authorizedBy, authorizedID)
	if err != nil {
		return nil, err
	} else if user == nil {
		return nil, nil
	}

	userRoles, err := u.UserRepo.FindAllUserRoles(user.UserID)
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
	userVO := entity.UserVO{
		AuthorizedBy: authorizedBy,
		AuthorizedID: authorizedID,
		Email:        email,
	}

	return u.UserRepo.SaveUser(userVO)
}
