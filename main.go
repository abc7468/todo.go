package main

import (
	"log"
	"net/http"

	"github.com/abc7468/todo.go/app"
)

func main() {
	m := app.MakeHandler("./test.db")
	defer m.Close()

	log.Print("Start")
	err := http.ListenAndServe(":3000", m)
	if err != nil {
		panic(err)
	}
}
