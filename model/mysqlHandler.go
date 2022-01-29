package model

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type sqlHandler struct {
	db *sql.DB
}

// dbHandler interface 구현
func (s *sqlHandler) addTodo(title string) *Todo {
	return nil
}

func (s *sqlHandler) getTodos() []*Todo {
	return nil
}

func (s *sqlHandler) removeTodo(id int) bool {
	return false
}

func (s *sqlHandler) completeTodo(id int, isComplete bool) bool {
	return false
}

func (s *sqlHandler) close() {
	s.db.Close()
}

func newMysqlHandler() dbHandler {
	db, err := sql.Open("mysql", "root:asdfasdf1!@tcp(127.0.0.1:3306)/study")
	if err != nil {
		log.Fatal("Unable to open!!!!!!!!!!!!!!!1")
	}
	results, err := db.Query("select * from category")
	if nil != err {
		log.Fatal("Error when fetching category table rows", err)
	}
	log.Print("results!!!!", results)
	// for results.Next() {
	// 쿼리해온 값의 next값이 없을 때 까지 실행

	// }
	// statement, _ := db.Prepare(`CREATE TABLE IF NOT EXISTS todos(
	// 	id INTEGER PRIMARY KEY AUTOINCREMENT,
	// 	title TEXT,
	// 	isComplete BOOLEAN,
	// 	createdAt DATETIME
	// )`)
	// statement.Exec()
	return &sqlHandler{db: db}
}
