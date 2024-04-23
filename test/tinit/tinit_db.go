package tinit

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"userService/usersvc/common/mysqldb"
	"userService/usersvc/common/vars"
	"userService/usersvc/models"
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

	v := reflect.ValueOf(models.TableNames)

	for i := 0; i < v.NumField(); i++ {
		_, err = db.Exec(fmt.Sprintf("TRUNCATE TABLE %v", v.Field(i).Interface()))
		checkErr(t, err)
	}
}
