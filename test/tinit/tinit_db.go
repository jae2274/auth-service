package tinit

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/jae2274/auth-service/auth_service/common/mysqldb"
	"github.com/jae2274/auth-service/auth_service/common/vars"
	"github.com/jae2274/auth-service/auth_service/models"
)

func DB(t *testing.T) *sql.DB {
	envVars, err := vars.Variables()
	CheckErr(t, err)

	sqlDB, err := mysqldb.DB(envVars.DbVars)
	CheckErr(t, err)

	ClearDB(t, sqlDB)

	return sqlDB
}

func CheckErr(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func ClearDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	CheckErr(t, err)

	defer func() {
		_, err := db.Exec("SET FOREIGN_KEY_CHECKS = 1")
		CheckErr(t, err)
	}()

	v := reflect.ValueOf(models.TableNames)

	for i := 0; i < v.NumField(); i++ {
		rawSql := fmt.Sprintf("TRUNCATE TABLE %v", v.Field(i).Interface())
		_, err = db.Exec(rawSql)
		CheckErr(t, err)
	}
}
