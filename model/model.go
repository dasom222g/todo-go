package model

import "time"

// inmemory 데이터를 struct에서 관리하며 메소드로 로직 정의
// var todoMap map[int]*Todo
// var currentID int

type Todo struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	IsComplete bool      `json:"is_complete"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// memoryHandler는 dbHandler 인터페이스를 구현함

type DBHandler interface {
	AddTodo(title string) *Todo
	GetTodos() []*Todo
	RemoveTodo(id int) bool
	CompleteTodo(id int, isComplete bool) bool
	Close()
}

// var handler dbHandler

func NewDBHandler(dbName string) DBHandler { // 패키지 실행 변수 초기화 이후로 실행됨
	// 같은 패키지에 있으므로 따로 import하지 않아도 인식됨
	return newMysqlHandler(dbName)
}
