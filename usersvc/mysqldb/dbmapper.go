package mysqldb

import (
	"database/sql"
	"userService/usersvc/domain"
	"userService/usersvc/entity"

	"github.com/jae2274/goutils/terr"
)

func SaveUser(tx *sql.Tx, user entity.UserVO) error {
	_, err := tx.Exec("INSERT INTO user (authorized_by, authorized_id, email) VALUES ( ?, ?, ?)", user.AuthorizedBy, user.AuthorizedID, user.Email)
	if err != nil {
		return terr.Wrap(err)
	}
	return nil //TODO
}

func FindByAuthorized(tx *sql.Tx, authorizedType domain.AuthorizedBy, authorizedID string) (*entity.UserVO, error) {
	row := tx.QueryRow(
		`SELECT user_id,
		authorized_by,
		authorized_id,
		email,
		create_date
	FROM user WHERE authorized_by = ? and authorized_id = ?`, authorizedType, authorizedID)

	var user entity.UserVO
	err := row.Scan(&user.UserID, &user.AuthorizedBy, &user.AuthorizedID, &user.Email, &user.CreateDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, terr.Wrap(err)
	}
	return &user, nil
}

func FindAllUserRoles(tx *sql.Tx, userID int64) ([]entity.UserRoleVO, error) {
	rows, err := tx.Query(`
	SELECT user_id,
    role_name,
    granted_type,
    granted_by,
    expiry_date
FROM user_role WHERE user_id=?`, userID)

	if err != nil {
		return nil, terr.Wrap(err)
	}
	defer rows.Close()

	var userRoles []entity.UserRoleVO
	for rows.Next() {
		var userRole entity.UserRoleVO
		err := rows.Scan(&userRole.UserID, &userRole.RoleName, &userRole.GrantedType, &userRole.GrantedBy, &userRole.ExpiryDate)
		if err != nil {
			return nil, terr.Wrap(err)
		}
		userRoles = append(userRoles, userRole)
	}

	return userRoles, nil //TODO
}
