package tinit

import (
	"database/sql"
	"testing"
	"userService/usersvc/mysqldb"
	"userService/usersvc/vars"
)

func DB(t *testing.T) *sql.DB {
	envVars, err := vars.Variables()
	checkErr(t, err)

	sqlDB, err := mysqldb.DB(envVars.DbVars)
	checkErr(t, err)

	ClearDB(t, sqlDB)

	return sqlDB
}

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func ClearDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	checkErr(t, err)

	defer func() {
		_, err := db.Exec("SET FOREIGN_KEY_CHECKS = 1")
		checkErr(t, err)
	}()

	_, err = db.Exec("TRUNCATE TABLE user_role")
	checkErr(t, err)

	_, err = db.Exec("TRUNCATE TABLE user")
	checkErr(t, err)
}
