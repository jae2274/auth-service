package mysqldb

import (
	"database/sql"
	"log"
	"strconv"
	"userService/usersvc/vars"
)

func DB(dbVars *vars.DBVars) (*sql.DB, error) {
	dataSourceName := dbVars.Username + ":" + dbVars.Password + "@tcp(" + dbVars.Host + ":" + strconv.FormatInt(dbVars.Port, 10) + ")/" + dbVars.Name
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}
