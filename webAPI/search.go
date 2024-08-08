package webAPI

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"text/template"

	"literary-lions-forum/databaseAPI"

	_ "github.com/mattn/go-sqlite3"
)

type Category struct {
	Name string
	Icon string
}

type Post struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Content  string `json:"content"`
	Username string `json:"username"`
}

func Search(w http.ResponseWriter, r *http.Request) {
	var payload HomePage
	var err error

	// Check if the user is logged in
	if isLoggedIn(r) {
		cookie, _ := r.Cookie("SESSION")
		payload.User = User{
			IsLoggedIn: true,
			Username:   databaseAPI.GetUser(database, cookie.Value),
		}
	} else {
		payload.User = User{IsLoggedIn: false}
	}

	// Get categories for the dropdown
	payload.Categories = databaseAPI.GetCategories(database)

	// Check if search parameters are present
	if r.URL.Query().Has("category") || r.URL.Query().Has("keywords") || r.URL.Query().Has("user") {
		// Perform the search query
		payload.Results, err = performSearch(database, r.URL.Query().Get("category"), r.URL.Query().Get("keywords"), r.URL.Query().Get("user"))
		if err != nil {
			http.Error(w, "Error performing search", http.StatusInternalServerError)
			log.Printf("Search error: %v", err) // Log the error
			return
		}

		// If the request is AJAX, return JSON
		if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
			w.Header().Set("Content-Type", "application/json")
			if len(payload.Results) > 0 {
				json.NewEncoder(w).Encode(payload.Results)
			} else {
				json.NewEncoder(w).Encode([]Post{}) // Return an empty array if no results
			}
			return
		}
	}

	// Parse and execute the template for regular requests
	t, err := template.ParseGlob("public/HTML/*.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Printf("Template error: %v", err) // Log the error
		return
	}
	t.ExecuteTemplate(w, "search.html", payload)
}

func performSearch(db *sql.DB, category, keywords, user string) ([]Post, error) {
	var results []Post
	var queryBuilder strings.Builder
	var args []interface{}

	queryBuilder.WriteString("SELECT id, title, categories, content, username FROM posts WHERE 1=1")

	if category != "" {
		queryBuilder.WriteString(" AND categories = ?")
		args = append(args, category)
	}

	if keywords != "" {
		queryBuilder.WriteString(" AND (title LIKE ? OR content LIKE ?)")
		args = append(args, "%"+keywords+"%", "%"+keywords+"%")
	}

	if user != "" {
		queryBuilder.WriteString(" AND username = ?")
		args = append(args, user)
	}

	rows, err := db.Query(queryBuilder.String(), args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Category, &post.Content, &post.Username); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		results = append(results, post)
	}

	return results, nil
}
