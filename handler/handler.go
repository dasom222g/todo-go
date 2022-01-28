package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/dasom222g/todo-go/check"
	"github.com/dasom222g/todo-go/model"
	"github.com/gorilla/pat"
	"github.com/unrolled/render"
)

var rd *render.Render

// var todoMap map[int]*Todo
// var currentID int

// type Todo struct {
// 	ID         int       `json:"id"`
// 	Title      string    `json:"title"`
// 	IsComplete bool      `json:"is_complete"`
// 	CreatedAt  time.Time `json:"created_at"`
// 	UpdatedAt  time.Time `json:"updated_at"`
// }

type Success struct {
	Success bool `json:"success"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func handleGetTodos(w http.ResponseWriter, r *http.Request) {
	todos := model.GetTodos()
	rd.JSON(w, http.StatusOK, todos)
}

func handleAddTodo(w http.ResponseWriter, r *http.Request) {
	todo := new(model.Todo)
	err := json.NewDecoder(r.Body).Decode(todo)
	if check.IsError(err, rd, w, http.StatusBadRequest) {
		return
	}
	newTodo := model.AddTodo(todo.Title)
	rd.JSON(w, http.StatusCreated, newTodo)
}

func handleRemoveTodo(w http.ResponseWriter, r *http.Request) {
	url := *r.URL // {    /todos/1  false %3Aid=1  }
	pathSlice := strings.Split(url.Path, "/")
	id, _ := strconv.Atoi(pathSlice[len(pathSlice)-1])
	ok := model.RemoveTodo(id)
	rd.JSON(w, http.StatusOK, &Success{ok})
}

func handleCompleteTodo(w http.ResponseWriter, r *http.Request) {
	// query메소드 값을 FormValue로 받음
	isComplete := r.FormValue("complete") == "true"
	url := *r.URL
	pathSlice := strings.Split(url.Path, "/")
	id, _ := strconv.Atoi(pathSlice[len(pathSlice)-1])
	ok := model.CompleteTodo(id, isComplete)
	rd.JSON(w, http.StatusOK, &Success{ok})
}

func NewHttpHandler() http.Handler {
	rd = render.New()
	// todoMap = make(map[int]*Todo)
	// currentID = 0

	mux := pat.New()
	mux.Get("/todos/{id:[0-9]+}", handleCompleteTodo)
	mux.Get("/todos", handleGetTodos)
	mux.Post("/todos", handleAddTodo)
	// mux.Delete("/todos{id}", handleRemoveTodo)
	mux.Delete("/todos/{id:[0-9]+}", handleRemoveTodo)
	mux.Get("/", indexHandler)
	return mux
}
