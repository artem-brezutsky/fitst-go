package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts []Article
var templates = []string{"templates/index.html", "templates/header.html", "templates/footer.html"}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(templates...)

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "index", nil)
}

func handleFunc() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}
