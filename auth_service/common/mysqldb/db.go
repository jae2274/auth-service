package mysqldb

import (
	"database/sql"
	"strconv"

	"github.com/jae2274/auth-service/auth_service/common/vars"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jae2274/goutils/terr"
)

func DB(dbVars *vars.DBVars) (*sql.DB, error) {
	dataSourceName := dbVars.Username + ":" + dbVars.Password + "@tcp(" + dbVars.Host + ":" + strconv.FormatInt(dbVars.Port, 10) + ")/" + dbVars.Name + "?parseTime=true"
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, terr.Wrap(err)
	}
	return db, nil
}
