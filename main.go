package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	"html/template"
	"log"
	"net/http"
	"os"
	_ "os"
)

func init() {
	// Load values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

type Article struct {
	Id                     uint16
	Title, Anons, FullText string
}

var posts []Article
var templates = []string{
	"templates/index.html",
	"templates/header.html",
	"templates/footer.html",
	"templates/create.html",
}

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(templates...)

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "index", nil)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles(templates...)

	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "create", nil)
}

func saveArticle(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	anons := r.FormValue("anons")
	fullText := r.FormValue("full_text")

	fmt.Println(title, anons, fullText)

	db := connectDb()

	insert, err := db.Query(fmt.Sprintf("INSERT INTO `articles` (`title`, `anons`, `full_text`) VALUES('%s', '%s', '%s')", title, anons, fullText))
	if err != nil {
		panic(err)
	}
	defer insert.Close()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleFunc() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/create", create)
	http.HandleFunc("/save_article", saveArticle)
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}

func connectDb() *sql.DB {
	dataSourceName, _ := os.LookupEnv("DB_SOURCE_NAME")
	dbDriver, _ := os.LookupEnv("DB_DRIVER")

	db, err := sql.Open(dbDriver, dataSourceName)
	if err != nil {
		panic(err)
	}

	return db
}
