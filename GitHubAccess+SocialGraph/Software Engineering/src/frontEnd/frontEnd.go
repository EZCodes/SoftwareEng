package main

import (
	"net/http"
	"log" 
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/pattern/", testHandler)
    log.Fatal(http.ListenAndServe(":8081", nil))
    
}