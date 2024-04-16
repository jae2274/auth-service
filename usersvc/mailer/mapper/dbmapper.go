package mapper

import (
	"bytes"
	"database/sql"
	"fmt"
	"userService/usersvc/common/entity"
)

func GetUserEMails(tx *sql.Tx, userIds []int64) ([]*entity.UserVO, error) {
	if len(userIds) == 0 {
		return []*entity.UserVO{}, nil
	}

	buff := bytes.NewBuffer([]byte{})
	for i, id := range userIds {
		if i > 0 {
			buff.WriteString(",")
		}
		buff.WriteString(fmt.Sprintf("%d", id))
	}
	rows, err := tx.Query("SELECT user_id, email FROM user WHERE user_id IN (?) AND agree_mail=?", buff.String(), true)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userEMails []*entity.UserVO
	for rows.Next() {
		var userEMail entity.UserVO
		err := rows.Scan(&userEMail.UserID, &userEMail.Email)
		if err != nil {
			return nil, err
		}
		userEMails = append(userEMails, &userEMail)
	}

	return userEMails, nil
}
