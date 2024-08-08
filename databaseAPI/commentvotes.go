package databaseAPI

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// HasCommentUpvoted checks if the user has upvoted a comment
func HasCommentUpvoted(database *sql.DB, username string, post_id int) bool {
	rows, err := database.Query("SELECT vote FROM comment_votes WHERE username = ? AND post_id = ? AND vote = 1", username, post_id)
	if err != nil {
		return false
	}
	defer rows.Close()

	vote := 0
	for rows.Next() {
		if err := rows.Scan(&vote); err != nil {
			return false
		}
	}

	return vote != 0
}

// HasCommentDownvoted checks if the user has downvoted a comment
func HasCommentDownvoted(database *sql.DB, username string, post_id int) bool {
	rows, err := database.Query("SELECT vote FROM comment_votes WHERE username = ? AND post_id = ? AND vote = -1", username, post_id)
	if err != nil {
		return false
	}
	defer rows.Close()

	vote := 0
	for rows.Next() {
		if err := rows.Scan(&vote); err != nil {
			return false
		}
	}

	return vote != 0
}

// RemoveCommentVote removes a vote from a comment
func RemoveCommentVote(database *sql.DB, post_id int, username string) {
	statement, err := database.Prepare("DELETE FROM comment_votes WHERE post_id = ? AND username = ?")
	if err != nil {
		return
	}
	defer statement.Close()

	_, err = statement.Exec(post_id, username)
	if err != nil {
		return
	}
}

// DecreaseCommentUpvotes decreases the upvotes of a comment by 1
func DecreaseCommentUpvotes(database *sql.DB, Id int) {
	statement, err := database.Prepare("UPDATE comments SET upvotes = upvotes - 1 WHERE Id = ?")
	if err != nil {
		return
	}
	defer statement.Close()

	_, err = statement.Exec(Id)
	if err != nil {
		return
	}
}

// DecreaseCommentDownvotes decreases the downvotes of a comment by 1
func DecreaseCommentDownvotes(database *sql.DB, Id int) {
	statement, err := database.Prepare("UPDATE comments SET downvotes = downvotes - 1 WHERE Id = ?")
	if err != nil {
		return
	}
	defer statement.Close()

	_, err = statement.Exec(Id)
	if err != nil {
		return
	}
}

// IncreaseCommentUpvotes increases the upvotes of a comment by 1
func IncreaseCommentUpvotes(database *sql.DB, Id int) {
	statement, err := database.Prepare("UPDATE comments SET upvotes = upvotes + 1 WHERE Id = ?")
	if err != nil {
		return
	}
	defer statement.Close()

	_, err = statement.Exec(Id)
	if err != nil {
		return
	}
}

// IncreaseCommentDownvotes increases the downvotes of a comment by 1
func IncreaseCommentDownvotes(database *sql.DB, Id int) {
	statement, err := database.Prepare("UPDATE comments SET downvotes = downvotes + 1 WHERE Id = ?")
	if err != nil {
		return
	}
	defer statement.Close()

	_, err = statement.Exec(Id)
	if err != nil {
		return
	}
}

// AddCommentVote adds a vote to the database for a comment
func AddCommentVote(database *sql.DB, post_id int, username string, vote int) {
	statement, err := database.Prepare("INSERT INTO comment_votes (username, post_id, vote) VALUES (?, ?, ?)")
	if err != nil {
		return
	}
	defer statement.Close()

	_, err = statement.Exec(username, post_id, vote)
	if err != nil {
		return
	}
}

// UpdateCommentVote updates the vote of a user for a comment
func UpdateCommentVote(database *sql.DB, post_id int, username string, vote int) {
	statement, err := database.Prepare("UPDATE comment_votes SET vote = ? WHERE post_id = ? AND username = ?")
	if err != nil {
		return
	}
	defer statement.Close()

	_, err = statement.Exec(vote, post_id, username)
	if err != nil {
		return
	}
}
