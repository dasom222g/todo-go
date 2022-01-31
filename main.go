package main

import (
	"net/http"

	"github.com/dasom222g/todo-go/handler"
)

var dbName string = "dasom"

func main() {
	// 처음 sqlHandler를 생성해준 main에서 close해줌
	m := handler.NewHttpHandler(dbName)
	defer m.Close()

	// n := negroni.Classic()
	// n.UseHandler(m)

	err := http.ListenAndServe(":3000", m)
	if err != nil {
		panic(err)
	}
}
