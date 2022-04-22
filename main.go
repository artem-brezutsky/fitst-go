package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
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
	checkErr(err)

	db := connectDb()
	res, err := db.Query("SELECT * FROM `articles`")
	checkErr(err)

	posts = []Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		checkErr(err)

		posts = append(posts, post)
	}

	defer db.Close()

	t.ExecuteTemplate(w, "index", posts)
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

func showPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "ID %v\n", vars["id"])
}

func handleFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", saveArticle).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", showPost).Methods("GET")

	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleFunc()
}

func connectDb() *sql.DB {
	dataSourceName, _ := os.LookupEnv("DB_SOURCE_NAME")
	dbDriver, _ := os.LookupEnv("DB_DRIVER")

	db, err := sql.Open(dbDriver, dataSourceName)
	checkErr(err)

	return db
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
