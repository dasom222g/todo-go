package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var contentType string = "application/json"

func addTodo(url, title string, assert *assert.Assertions) *Todo {
	s := fmt.Sprintf(`{"title": "%s!!", "is_complete": false}`, title)
	res, err := http.Post(url, contentType, strings.NewReader(s))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)
	newTodo := new(Todo)
	err = json.NewDecoder(res.Body).Decode(&newTodo)
	assert.NoError(err)
	assert.Equal(title+"!!", newTodo.Title)
	return newTodo
}

func TestHandleAddTodo(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(NewHttpHandler())
	defer ts.Close()

	newTodo1 := addTodo(ts.URL+"/todos", "Buy some milk", assert)
	newTodo2 := addTodo(ts.URL+"/todos", "Go to sleep", assert)
	res, err := http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	// 추가한 데이터들이 잘 들어가 있는지 확인
	todos := []*Todo{}
	err = json.NewDecoder(res.Body).Decode(&todos)
	assert.NoError(err)
	assert.Equal(2, len(todos))
	for _, todo := range todos {
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
	ts := httptest.NewServer(NewHttpHandler())
	defer ts.Close()

	newTodo1 := addTodo(ts.URL+"/todos", "Buy some milk", assert)
	newTodo2 := addTodo(ts.URL+"/todos", "Go to sleep", assert)

	res, err := http.Get(ts.URL + "/todos/" + strconv.Itoa(newTodo1.ID) + "?complete=true")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	res, err = http.Get(ts.URL + "/todos/" + strconv.Itoa(newTodo2.ID) + "?complete=false")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	todos := []*Todo{}
	res, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(&todos)
	assert.NoError(err)
	for _, todo := range todos {
		if todo.ID == newTodo1.ID {
			assert.True(todo.IsComplete)
		} else if todo.ID == newTodo2.ID {
			assert.False(todo.IsComplete)
		}
	}

}

func TestHandleRemoveTodo(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(NewHttpHandler())
	defer ts.Close()
	newTodo1 := addTodo(ts.URL+"/todos", "Buy some milk", assert)
	newTodo2 := addTodo(ts.URL+"/todos", "Go to sleep", assert)

	req, _ := http.NewRequest("DELETE", ts.URL+"/todos/"+strconv.Itoa(newTodo1.ID), nil)
	res, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)

	// 데이터 삭제후 todos get해서 값이 없는 지 확인
	todos := []*Todo{}
	res, err = http.Get(ts.URL + "/todos")
	assert.NoError(err)
	assert.Equal(http.StatusOK, res.StatusCode)
	err = json.NewDecoder(res.Body).Decode(&todos)
	assert.NoError(err)
	for _, todo := range todos {
		assert.Equal(newTodo2.ID, todo.ID)
	}

}
