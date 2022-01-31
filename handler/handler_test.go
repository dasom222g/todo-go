package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/dasom222g/todo-go/model"
	"github.com/stretchr/testify/assert"
)

var contentType string = "application/json"
var dbName string = "dasom"

func addTodo(url, title string, assert *assert.Assertions) *model.Todo {
	getSessionId = func(r *http.Request) string {
		return "test_session_id"
	}
	s := fmt.Sprintf(`{"title": "%s!!", "is_complete": false}`, title)
	res, err := http.Post(url, contentType, strings.NewReader(s))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)
	newTodo := new(model.Todo)
	err = json.NewDecoder(res.Body).Decode(newTodo)
	assert.NoError(err)
	assert.Equal(title+"!!", newTodo.Title)
	return newTodo
}

func TestHandleAddTodo(t *testing.T) {
	assert := assert.New(t)
	h := NewHttpHandler(dbName)
	defer h.Close()
	ts := httptest.NewServer(h)
	defer ts.Close()

	newTodo1 := addTodo(ts.URL+"/todos", "Buy some milk", assert)
	newTodo2 := addTodo(ts.URL+"/todos", "Go to sleep", assert)
	res, err := http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	// 추가한 데이터들이 잘 들어가 있는지 확인
	todos := []*model.Todo{}
	err = json.NewDecoder(res.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(2, len(todos))
	for _, todo := range todos {
		fmt.Println(todo)
		if todo.ID == newTodo1.ID {
			assert.Equal("Buy some milk!!", todo.Title)
		} else if todo.ID == newTodo2.ID {
			assert.Equal("Go to sleep!!", todo.Title)
		} else {
			assert.Error(fmt.Errorf("testID should be id1 or id2"))
		}
	}
}

func TestHandleCompleteTodo(t *testing.T) {
	assert := assert.New(t)
	h := NewHttpHandler(dbName)
	defer h.Close()
	ts := httptest.NewServer(h)
	defer ts.Close()

	newTodo3 := addTodo(ts.URL+"/todos", "Buy some milk tt", assert)
	newTodo4 := addTodo(ts.URL+"/todos", "Go to sleep tt", assert)

	res, err := http.Get(ts.URL + "/todos/" + strconv.Itoa(newTodo3.ID) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	res, err = http.Get(ts.URL + "/todos/" + strconv.Itoa(newTodo4.ID) + "?complete=false")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	res, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	todos := []*model.Todo{}
	err = json.NewDecoder(res.Body).Decode(&todos)
	assert.NoError(err)
	for _, todo := range todos {
		if todo.ID == newTodo3.ID {
			assert.True(todo.IsComplete)
		} else if todo.ID == newTodo4.ID {
			assert.False(todo.IsComplete)
		}
	}

}

func TestHandleRemoveTodo(t *testing.T) {
	assert := assert.New(t)
	h := NewHttpHandler(dbName)
	defer h.Close()
	ts := httptest.NewServer(h)
	defer ts.Close()
	newTodo5 := addTodo(ts.URL+"/todos", "Buy some milk454", assert)
	newTodo6 := addTodo(ts.URL+"/todos", "Go to sleep", assert)

	req, _ := http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(newTodo5.ID), nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	success := new(Success)
	err = json.NewDecoder(res.Body).Decode(success)
	assert.NoError(err)

	// 데이터 삭제후 todos get해서 값이 없는 지 확인
	res, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	todos := []*model.Todo{}
	err = json.NewDecoder(res.Body).Decode(&todos)
	assert.NoError(err)
	for _, todo := range todos {
		// 총 5번 돔..
		fmt.Println("todo!!", todo)
		if todo.ID == newTodo6.ID {
			assert.Equal("Go to sleep!!", todo.Title)
		} else if todo.ID == newTodo5.ID {
			assert.Error(fmt.Errorf("Not exists ID"))
		}
	}

}
