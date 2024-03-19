package mysqldb

import (
	"database/sql"
	"userService/usersvc/domain"
	"userService/usersvc/entity"
)

func SaveUser(tx *sql.Tx, user entity.UserVO) error {
	return nil //TODO
}

func FindByAuthorized(tx *sql.Tx, authorizedType domain.AuthorizedBy, authorizedID string) (*entity.UserVO, error) {
	return nil, nil //TODO
}

func FindAllUserRoles(tx *sql.Tx, userID int64) ([]entity.UserRoleVO, error) {
	return nil, nil //TODO
}
