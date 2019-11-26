package main

import (
	"html/template"
	"net/http"
	"fmt"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("src/frontEnd/index.html")
	if err != nil {
		fmt.Println("Error parsing template")
	}
	t.Execute(w, nil)
}