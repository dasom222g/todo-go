package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var contentType string = "application/json"

func TestHandleAddTodo(t *testing.T) {
	assert := assert.New(t)
	ts := httptest.NewServer(NewHttpHandler())
	defer ts.Close()

	s := `{"title": "Buy some milk", "is_complete": false}`
	res, err := http.Post(ts.URL+"/todos", contentType, strings.NewReader(s))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, res.StatusCode)
	newTodo := new(Todo)
	json.NewDecoder(res.Body).Decode(newTodo)
	assert.Equal("Buy some milk", newTodo.Title)
}
