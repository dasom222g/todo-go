package handler

import (
	"net/http"

	"github.com/gorilla/pat"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/todo.html", http.StatusTemporaryRedirect)
}

func NewHttpHandler() http.Handler {
	mux := pat.New()
	mux.Get("/", indexHandler)
	return mux
}
