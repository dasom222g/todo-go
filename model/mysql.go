package model

import "database/sql"

type sqlHandler struct {
	db *sql.DB
}

func newMysqlHandler() *sqlHandler {
	db, err := sql.Open("mysql", "root:asdfasdf1!@tcp(127.0.0.1:3306)/test")
	// defer db.Close()
	if err != nil {
		panic(err)
	}
	return &sqlHandler{db}
}
