package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dasom222g/todo-go/check"
	_ "github.com/go-sql-driver/mysql"
)

type sqlHandler struct {
	db *sql.DB
}

// dbHandler interface 구현
func (s *sqlHandler) AddTodo(title string) *Todo {
	results, err := s.db.Exec("INSERT INTO todos (title, is_complete) VALUES (?, ?)", title, false)
	check.CheckError(err)
	id, err := results.LastInsertId()
	check.CheckError(err)
	todo := &Todo{}
	todo.ID = int(id)
	todo.Title = title
	todo.IsComplete = false
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	return todo
}

func (s *sqlHandler) GetTodos() []*Todo {
	rows, err := s.db.Query("SELECT id, title, is_complete, created_at, updated_at FROM todos")
	check.CheckError(err)

	todos := []*Todo{}
	for rows.Next() {
		// rows를 한줄 씩 읽으며 다음 row가 있는지 체크하고 있을 때 까지 실행
		var todo Todo
		// 각 row의 데이터들을 todo에 넣어줌
		err := rows.Scan(&todo.ID, &todo.Title, &todo.IsComplete, &todo.CreatedAt, &todo.UpdatedAt)
		check.CheckError(err)
		todos = append(todos, &todo)
	}
	return todos
}

func (s *sqlHandler) RemoveTodo(id int) bool {
	statement, err := s.db.Prepare("DELETE FROM todos WHERE id=?")
	check.CheckError(err)
	results, err := statement.Exec(id)
	check.CheckError(err)
	changeRowCount, err := results.RowsAffected() // 변경된 row의 갯수
	check.CheckError(err)
	return changeRowCount == 1
}

func (s *sqlHandler) CompleteTodo(id int, isComplete bool) bool {
	statement, err := s.db.Prepare("UPDATE todos SET is_complete=? WHERE id=?")
	check.CheckError(err)
	results, err := statement.Exec(isComplete, id)
	check.CheckError(err)
	changeRowCount, err := results.RowsAffected() // 변경된 row의 갯수
	return changeRowCount == 1
}

func (s *sqlHandler) Close() {
	s.db.Close()
}

func newMysqlHandler(dbName string) DBHandler {
	db, err := sql.Open("mysql", fmt.Sprintf("root:asdfasdf1!@tcp(127.0.0.1:3306)/%s?parseTime=true", dbName))
	check.CheckError(err)
	// create table
	query := `CREATE TABLE IF NOT EXISTS todos (
		id BIGINT PRIMARY KEY AUTO_INCREMENT,
		title VARCHAR(255),
		is_complete BOOLEAN,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`
	if _, err := db.Exec(query); err != nil {
		panic(err)
	}
	return &sqlHandler{db: db}
}
