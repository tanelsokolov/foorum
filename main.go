package main

import (
	"database/sql"
	"fmt"
	"literary-lions-forum/databaseAPI"
	"literary-lions-forum/webAPI"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Post struct {
	Id         int
	Username   string
	Title      string
	Categories []string
	Content    string
	CreatedAt  string
	UpVotes    int
	DownVotes  int
	Comments   []Comment
}

type Comment struct {
	Id        int
	PostId    int
	Username  string
	Content   string
	CreatedAt string
	UpVotes   int
	DownVotes int
}

// Database
var database *sql.DB

func main() {
	// check if DB exists
	var _, err = os.Stat("database.db")

	// create DB if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create("database.db")
		if err != nil {
			return
		}
		defer file.Close()
	}

	database, _ = sql.Open("sqlite3", "./database.db")

	databaseAPI.CreateUsersTable(database)
	databaseAPI.CreatePostTable(database)
	databaseAPI.CreateCommentTable(database)
	databaseAPI.CreateCommentVoteTable(database)
	databaseAPI.CreateVoteTable(database)
	databaseAPI.CreateCategoriesTable(database)
	databaseAPI.CreateCategories(database)
	databaseAPI.CreateCategoriesIcons(database)

	webAPI.SetDatabase(database)

	fs := http.FileServer(http.Dir("public"))
	router := http.NewServeMux()
	fmt.Println("Starting server on port 8000")

	router.HandleFunc("/", webAPI.Index)
	router.HandleFunc("/register", webAPI.Register)
	router.HandleFunc("/login", webAPI.Login)
	router.HandleFunc("/post", webAPI.DisplayPost)
	router.HandleFunc("/filter", webAPI.GetPostsByApi)
	router.HandleFunc("/newpost", webAPI.NewPost)
	router.HandleFunc("/search", webAPI.Search)
	router.HandleFunc("/api/register", webAPI.RegisterApi)
	router.HandleFunc("/api/login", webAPI.LoginApi)
	router.HandleFunc("/api/logout", webAPI.LogoutAPI)
	router.HandleFunc("/api/createpost", webAPI.CreatePostApi)
	router.HandleFunc("/api/comments", webAPI.CommentsApi)
	router.HandleFunc("/api/comments/vote", webAPI.CommentsVoteApi)
	router.HandleFunc("/api/vote", webAPI.VoteApi)
	router.HandleFunc("/api/deletepost", webAPI.DeletePostApi)

	router.Handle("/public/", http.StripPrefix("/public/", fs))
	http.ListenAndServe(":8000", router)
}
