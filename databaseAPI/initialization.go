package databaseAPI

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// CreateUsersTable creates the users table
func CreateUsersTable(database *sql.DB) {
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, email TEXT, password TEXT, cookie TEXT, expires TEXT)")
	statement.Exec()
}

// CreatePostTable create post table
func CreatePostTable(database *sql.DB) {
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS posts (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, title TEXT, categories TEXT, content TEXT, created_at TEXT, upvotes INTEGER, downvotes INTEGER)")
	statement.Exec()
}

// CreateCommentTable creates a comment table
func CreateCommentTable(database *sql.DB) {
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS comments (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, post_id INTEGER, content TEXT, created_at TEXT, upvotes INTEGER, downvotes INTEGER)")
	statement.Exec()
}

// CreateVoteTable create the vote table into given database
func CreateVoteTable(database *sql.DB) {
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS votes (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, post_id INTEGER, vote INTEGER)")
	statement.Exec()
}

// CreateCommentVoteTable create the vote table into given database
func CreateCommentVoteTable(database *sql.DB) {
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS comment_votes (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, post_id INTEGER, vote INTEGER)")
	statement.Exec()
}

// CreateCategoriesTable create the categories' table into given database
func CreateCategoriesTable(database *sql.DB) {
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS categories (id INTEGER PRIMARY KEY, name TEXT, icon TEXT)")
	statement.Exec()
}

// CreateCategories creates categories in the database
func CreateCategories(database *sql.DB) {
	statement, _ := database.Prepare("INSERT INTO categories (name) SELECT ? WHERE NOT EXISTS (SELECT 1 FROM categories WHERE name = ?)")
	statement.Exec("General", "General")
	statement.Exec("Technology", "Technology")
	statement.Exec("Science", "Science")
	statement.Exec("Sports", "Sports")
	statement.Exec("Gaming", "Gaming")
	statement.Exec("Music", "Music")
}

// createCategoriesIcons creates categories' icons in the database
func CreateCategoriesIcons(database *sql.DB) {
	statement, _ := database.Prepare("UPDATE categories SET icon = ? WHERE name = ?")
	statement.Exec("fa-globe", "General")
	statement.Exec("fa-laptop", "Technology")
	statement.Exec("fa-flask", "Science")
	statement.Exec("fa-futbol-o", "Sports")
	statement.Exec("fa-gamepad", "Gaming")
	statement.Exec("fa-music", "Music")
}
