package model

import "database/sql"

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
	db, err := sql.Open("mysql", "root:asdfasdf1!@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}
	return &sqlHandler{db: db}
}
