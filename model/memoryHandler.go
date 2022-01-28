package model

import "time"

type memoryHandler struct {
	todoMap   map[int]*Todo
	currentID int
}

// interface 구현
func (m *memoryHandler) addTodo(title string) *Todo {
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

func (m *memoryHandler) getTodos() []*Todo {
	todos := []*Todo{}
	if len(m.todoMap) == 0 {
		return todos
	}
	for _, value := range m.todoMap {
		todos = append(todos, value)
	}
	return todos
}

func (m *memoryHandler) removeTodo(id int) bool {
	if _, exists := m.todoMap[id]; exists {
		delete(m.todoMap, id)
		return true
	}
	return false
}

func (m *memoryHandler) completeTodo(id int, isComplete bool) bool {
	if todo, exists := m.todoMap[id]; exists {
		todo.IsComplete = isComplete
		return true
	}
	return false
}

// 생성 될때 초기화
func newMemoryHandler() dbHandler {
	m := &memoryHandler{}
	m.currentID = 0
	m.todoMap = make(map[int]*Todo)
	return m
}
