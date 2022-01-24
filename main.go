package main

import (
	"net/http"

	"github.com/dasom222g/todo-go/handler"
	"github.com/urfave/negroni"
)

func main() {
	mux := handler.NewHttpHandler()
	n := negroni.Classic()
	n.UseHandler(mux)

	err := http.ListenAndServe(":3000", n)
	if err != nil {
		panic(err)
	}
}
