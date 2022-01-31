package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dasom222g/todo-go/check"
	"github.com/dasom222g/todo-go/model"
	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
)

var rd *render.Render = render.New()
var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

type Handler struct {
	http.Handler // 임베디드 형태
	DB           model.DBHandler
}

type Success struct {
	Success bool `json:"success"`
}

func (h *Handler) indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func (h *Handler) handleGetTodos(w http.ResponseWriter, r *http.Request) {
	todos := h.DB.GetTodos()
	rd.JSON(w, http.StatusOK, todos)
}

func (h *Handler) handleAddTodo(w http.ResponseWriter, r *http.Request) {
	todo := new(model.Todo)
	err := json.NewDecoder(r.Body).Decode(todo)
	if check.IsError(err, rd, w, http.StatusBadRequest) {
		return
	}
	newTodo := h.DB.AddTodo(todo.Title)
	rd.JSON(w, http.StatusCreated, newTodo)
}

func (h *Handler) handleRemoveTodo(w http.ResponseWriter, r *http.Request) {
	url := *r.URL // {    /todos/1  false %3Aid=1  }
	pathSlice := strings.Split(url.Path, "/")
	id, _ := strconv.Atoi(pathSlice[len(pathSlice)-1])
	ok := h.DB.RemoveTodo(id)
	rd.JSON(w, http.StatusOK, &Success{ok})
}

func (h *Handler) handleCompleteTodo(w http.ResponseWriter, r *http.Request) {
	// query메소드 값을 FormValue로 받음
	isComplete := r.FormValue("complete") == "true"
	url := *r.URL
	pathSlice := strings.Split(url.Path, "/")
	id, _ := strconv.Atoi(pathSlice[len(pathSlice)-1])
	ok := h.DB.CompleteTodo(id, isComplete)
	rd.JSON(w, http.StatusOK, &Success{ok})
}

func (h *Handler) Close() {
	h.DB.Close()
}

func getSessionId(r *http.Request) string {
	session, err := store.Get(r, "session")
	if err != nil {
		return ""
	}
	val := session.Values["id"]
	if val == nil {
		return ""
	}
	// 데이터형식 모를경우 .(string) 으로 명시
	return val.(string)
}

func NewHttpHandler(dbName string) *Handler {
	mux := pat.New()
	// 초기화
	h := &Handler{
		Handler: mux,
		DB:      model.NewDBHandler(dbName),
	}
	mux.Get("/todos/{id:[0-9]+}", h.handleCompleteTodo)
	mux.Get("/todos", h.handleGetTodos)
	mux.Post("/todos", h.handleAddTodo)
	mux.Delete("/todos/{id:[0-9]+}", h.handleRemoveTodo)
	mux.HandleFunc("/auth/google/login", handleGoogleLogin)
	mux.HandleFunc("/auth/google/callback", handleGoogleLoginCallback)
	mux.Get("/", h.indexHandler)
	return h
}
