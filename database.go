package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "gopkg.in/goracle.v2"
)

var (
	dbusername = "*******"
	dbpassword = "*******"
	dbsid      = "*******"
)

/*
ConnectDB connects to the database and returns a sql.DB pointer
*/
func ConnectDB() *sql.DB {

	connString := fmt.Sprintf("%s/%s/@//****:****/%s",
		dbusername,
		dbpassword,
		dbsid)

	db, err := sql.Open("goracle", connString)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
