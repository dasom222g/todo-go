package model

import "time"

var todoMap map[int]*Todo
var currentID int

type Todo struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	IsComplete bool      `json:"is_complete"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func init() {
	// 패키지 실행 변수 초기화 이후로 실행됨
	todoMap = make(map[int]*Todo)
	currentID = 0
}

func AddTodo(title string) *Todo {
	todo := &Todo{}

	currentID++
	todo.ID = currentID
	todo.Title = title
	todo.IsComplete = false
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	todoMap[todo.ID] = todo
	return todo
}

func GetTodos() []*Todo {
	todos := []*Todo{}
	if len(todoMap) == 0 {
		return todos
	}
	for _, value := range todoMap {
		todos = append(todos, value)
	}
	return todos
}

func RemoveTodo(id int) bool {
	if _, exists := todoMap[id]; exists {
		delete(todoMap, id)
		return true
	}
	return false
}

func CompleteTodo(id int, isComplete bool) bool {
	if todo, exists := todoMap[id]; exists {
		todo.IsComplete = isComplete
		return true
	}
	return false
}
