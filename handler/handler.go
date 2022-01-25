package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dasom222g/todo-go/check"
	"github.com/gorilla/pat"
	"github.com/unrolled/render"
)

var rd *render.Render
var todoMap map[int]*Todo
var currentID int

type Todo struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	IsComplete bool      `json:"is_complete"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Success struct {
	Success bool `json:"success"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func handleGetTodos(w http.ResponseWriter, r *http.Request) {
	todos := []*Todo{}
	if len(todoMap) == 0 {
		rd.JSON(w, http.StatusOK, todos)
		return
	}

	for _, value := range todoMap {
		todos = append(todos, value)
	}
	rd.JSON(w, http.StatusOK, todos)
}

func handleAddTodo(w http.ResponseWriter, r *http.Request) {
	todo := new(Todo)
	err := json.NewDecoder(r.Body).Decode(todo)
	if check.IsError(err, rd, w, http.StatusBadRequest) {
		return
	}
	currentID++
	todo.ID = currentID
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	todoMap[todo.ID] = todo
	rd.JSON(w, http.StatusOK, todo)
}

func handleRemoveTodo(w http.ResponseWriter, r *http.Request) {
	url := *r.URL
	pathSlice := strings.Split(url.Path, "/")
	id, _ := strconv.Atoi(pathSlice[len(pathSlice)-1])
	log.Print("url: ", url)
	log.Printf("id: %d", id)
	if _, exists := todoMap[id]; exists {
		// 해당 아이템이 있는 경우
		delete(todoMap, id)
		rd.JSON(w, http.StatusOK, &Success{true})
	} else {
		rd.JSON(w, http.StatusOK, &Success{false})
	}
}

func NewHttpHandler() http.Handler {
	rd = render.New()
	todoMap = make(map[int]*Todo)
	currentID = 0

	mux := pat.New()
	mux.Get("/todos", handleGetTodos)
	mux.Post("/todos", handleAddTodo)
	// mux.Delete("/todos{id}", handleRemoveTodo)
	mux.Delete("/todos/{id:[0-9]+}", handleRemoveTodo)
	mux.Get("/", indexHandler)
	return mux
}
