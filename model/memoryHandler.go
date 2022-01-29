package model

import "time"

type memoryHandler struct {
	todoMap   map[int]*Todo
	currentID int
}

// dbHandler interface 구현
func (m *memoryHandler) AddTodo(title string) *Todo {
	todo := &Todo{}

	m.currentID++
	id := m.currentID
	todo.ID = id
	todo.Title = title
	todo.IsComplete = false
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	m.todoMap[id] = todo
	return todo
}

func (m *memoryHandler) GetTodos() []*Todo {
	todos := []*Todo{}
	if len(m.todoMap) == 0 {
		return todos
	}
	for _, value := range m.todoMap {
		todos = append(todos, value)
	}
	return todos
}

func (m *memoryHandler) RemoveTodo(id int) bool {
	if _, exists := m.todoMap[id]; exists {
		delete(m.todoMap, id)
		return true
	}
	return false
}

func (m *memoryHandler) CompleteTodo(id int, isComplete bool) bool {
	if todo, exists := m.todoMap[id]; exists {
		todo.IsComplete = isComplete
		return true
	}
	return false
}

func (m *memoryHandler) Close() {}

// 생성 될때 초기화
func newMemoryHandler() DBHandler {
	m := &memoryHandler{}
	m.currentID = 0
	m.todoMap = make(map[int]*Todo)
	return m
}
