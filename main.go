package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	PORT     = ":8080"
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "bobbyshmurda66"
	dbname   = "blog"
)

var database *sql.DB

func serveBlog(w http.ResponseWriter, r *http.Request) {

	parameters := mux.Vars(r)
	blogId := parameters["id"]
	currentBlog := Blog{}

	sqlStatement := `SELECT id,image,title,short_text,long_text,date FROM blog WHERE id=$1`
	row := database.QueryRow(sqlStatement, blogId)
	err := row.Scan(&currentBlog.Id, &currentBlog.Image, &currentBlog.Title, &currentBlog.ShortText, &currentBlog.LongText, &currentBlog.Date)

	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println("Couldn't get page !")

	}

	sqlStatementComment := `SELECT name, email, text FROM comment WHERE blog_id=$1`
	comments, err := database.Query(sqlStatementComment, currentBlog.Id)
	if err != nil {
		log.Println(err)
	}

	for comments.Next() {
		var comment Comment
		comments.Scan(&comment.Name, &comment.Email, &comment.Text)
		currentBlog.Comments = append(currentBlog.Comments, comment)
	}

	temp, _ := template.ParseFiles("templates/blog.html")
	temp.Execute(w, currentBlog)

}
func redirectHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", 301)
}

func homePage(w http.ResponseWriter, r *http.Request) {
	allBlogs := []Blog{}
	sqlStatement := `SELECT id,image,title,short_text,long_text,date FROM blog`
	rows, err := database.Query(sqlStatement)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {

		var blog Blog
		rows.Scan(&blog.Id, &blog.Image, &blog.Title, &blog.ShortText, &blog.LongText, &blog.Date)
		allBlogs = append(allBlogs, blog)

	}

	temp, _ := template.ParseFiles("templates/home.html")
	temp.Execute(w, allBlogs)

}

func ApiCommentPost(w http.ResponseWriter, r *http.Request) {
	var commentAdded bool
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	blog_id := r.FormValue("id")
	name := r.FormValue("name")
	email := r.FormValue("email")
	comment := r.FormValue("comment")
	date := time.Now()

	LastInsertId := 0
	sqlStatement := "INSERT INTO comment (blog_id,email,text,name,date) VALUES ($1,$2,$3,$4,$5) RETURNING id"
	row := database.QueryRow(sqlStatement, blog_id, email, comment, name, date).Scan(&LastInsertId)

	if row != nil {
		log.Println(err.Error())
	}
	id := LastInsertId
	if id == 0 {
		commentAdded = false
	} else {
		commentAdded = true
	}

	commentAddedBool := strconv.FormatBool(commentAdded)
	resp := make(map[string]string)
	resp["id"] = strconv.Itoa(id)
	resp["added"] = commentAddedBool
	JsonPrintResponse, _ := json.Marshal(resp)
	w.Header().Set("Content-type", "application/json")
	fmt.Fprintln(w, string(JsonPrintResponse))
}

func main() {

	databaseConnection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", databaseConnection)
	if err != nil {
		log.Fatal(err)
	}
	database = db

	router := mux.NewRouter()
	router.HandleFunc("/blog/{id:[0-9]+}", serveBlog)
	router.HandleFunc("/home", homePage)
	router.HandleFunc("/", redirectHome)

	router.HandleFunc("/api/comments", ApiCommentPost).Methods("POST")
	http.Handle("/", router)

	http.ListenAndServe(PORT, nil)
}
