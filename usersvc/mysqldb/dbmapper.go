package mysqldb

import (
	"database/sql"
	"userService/usersvc/domain"
	"userService/usersvc/entity"
)

type DBMapper interface {
	SaveUser(user entity.UserVO) error
	FindByAuthorized(authorizedType domain.AuthorizedBy, authorizedID string) (*entity.UserVO, error)
	FindAllUserRoles(userID int64) ([]entity.UserRoleVO, error)
}

type DBMapperImpl struct {
	db *sql.DB
}

func NewDBMapper(db *sql.DB) DBMapper {
	return &DBMapperImpl{
		db: db,
	}
}

func (d *DBMapperImpl) SaveUser(user entity.UserVO) error {
	return nil //TODO
}

func (d *DBMapperImpl) FindByAuthorized(authorizedType domain.AuthorizedBy, authorizedID string) (*entity.UserVO, error) {
	return nil, nil //TODO
}

func (d *DBMapperImpl) FindAllUserRoles(userID int64) ([]entity.UserRoleVO, error) {
	return nil, nil //TODO
}
