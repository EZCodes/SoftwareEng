package main

import (
	"html/template"
	"net/http"
	"fmt"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello World")
}

func viewHandler(w http.ResponseWriter, r *http.Request) {

    fmt.Fprintf(w, "<h1>This is</h1><div>View Page</div>")
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("src/source/index.html")
	if err != nil {
		fmt.Println("Error parsing template")
	}
	t.Execute(w, nil)
}