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
var showPost = Article{}
var templates = []string{
	"templates/header.html",
	"templates/footer.html",
}

func index(w http.ResponseWriter, r *http.Request) {
	templates = append(templates, "templates/index.html")
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

	var i int = 10
	for ; i < 100; i++ {
		switch i {
		case 0:
			fmt.Println(fmt.Sprintf("Zero: %d", i))
		case 1:
			fmt.Println(fmt.Sprintf("One: %d", i))
		case 2:
			fmt.Println(fmt.Sprintf("Two: %d", i))
		case 3:
			fmt.Println(fmt.Sprintf("Three: %d", i))
		case 4:
			fmt.Println(fmt.Sprintf("Four: %d", i))
		case 5:
			fmt.Println(fmt.Sprintf("Five: %d", i))
		default:
			fmt.Println("Unknown Number")
		}
	}

	defer db.Close()

	t.ExecuteTemplate(w, "index", posts)
}

func create(w http.ResponseWriter, r *http.Request) {
	templates = append(templates, "templates/create.html")
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

func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templates = append(templates, "templates/show.html")
	t, err := template.ParseFiles(templates...)
	checkErr(err)

	db := connectDb()
	defer db.Close()

	res, err := db.Query(fmt.Sprintf("SELECT * FROM `articles` WHERE `id` = %s", vars["id"]))
	checkErr(err)

	showPost = Article{}
	for res.Next() {
		var post Article
		err = res.Scan(&post.Id, &post.Title, &post.Anons, &post.FullText)
		checkErr(err)

		showPost = post
	}

	t.ExecuteTemplate(w, "show", showPost)
}

func handleFunc() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index).Methods("GET")
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", saveArticle).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

	http.Handle("/", rtr)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	http.ListenAndServe(":8080", nil)
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

func main() {
	handleFunc()
}
