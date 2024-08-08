package webAPI

import (
	"fmt"
	"literary-lions-forum/databaseAPI"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Vote struct {
	PostId int
	Vote   int
}

type CommentVote struct {
	PostId int
	Vote   int
}

// CreatePostApi creates a post
func CreatePostApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	cookie, _ := r.Cookie("SESSION")
	username := databaseAPI.GetUser(database, cookie.Value)
	title := r.FormValue("title")
	content := r.FormValue("content")
	categories := r.Form["categories[]"]
	validCategories := databaseAPI.GetCategories(database)
	for _, category := range categories {
		// if string not in array, return error
		if !inArray(category, validCategories) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Invalid category : " + category))
			return
		}
	}
	stringCategories := strings.Join(categories, ",")
	now := time.Now()
	databaseAPI.CreatePost(database, username, title, stringCategories, content, now)
	fmt.Println("Post created by " + username + " with title " + title + " at " + now.Format("2006-01-02 15:04:05"))
	http.Redirect(w, r, "/filter?by=myposts", http.StatusFound)
}

// CommentsApi creates a comment
func CommentsApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("ParseForm() error: %v", err), http.StatusInternalServerError)
		return
	}
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	cookie, _ := r.Cookie("SESSION")
	username := databaseAPI.GetUser(database, cookie.Value)
	postId := r.FormValue("postId")
	content := r.FormValue("content")
	now := time.Now()
	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid postId: %v", err), http.StatusBadRequest)
		return
	}

	Id, err := databaseAPI.AddComment(database, username, postIdInt, content, now)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error adding comment: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Comment created by %s with ID %d on post %d at %s\n", username, Id, postIdInt, now.Format("2006-01-02 15:04:05"))
	http.Redirect(w, r, "/post?id="+postId, http.StatusFound)
}

// VoteApi api to vote on a post
func VoteApi(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !isLoggedIn(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	cookie, _ := r.Cookie("SESSION")
	username := databaseAPI.GetUser(database, cookie.Value)
	postId := r.FormValue("postId")
	postIdInt, _ := strconv.Atoi(postId)
	vote := r.FormValue("vote")
	voteInt, _ := strconv.Atoi(vote)
	now := time.Now().Format("2006-01-02 15:04:05")

	switch voteInt {
	case 1:
		if databaseAPI.HasUpvoted(database, username, postIdInt) {
			databaseAPI.RemoveVote(database, postIdInt, username)
			databaseAPI.DecreaseUpvotes(database, postIdInt)
			fmt.Println("Removed upvote from " + username + " on post " + postId + " at " + now)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Vote removed"))
			return
		}
		if databaseAPI.HasDownvoted(database, username, postIdInt) {
			databaseAPI.DecreaseDownvotes(database, postIdInt)
			databaseAPI.IncreaseUpvotes(database, postIdInt)
			databaseAPI.UpdateVote(database, postIdInt, username, 1)
			fmt.Println(username + " upvoted" + " on post " + postId + " at " + now)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Upvote added"))
			return
		}
		databaseAPI.IncreaseUpvotes(database, postIdInt)
		databaseAPI.AddVote(database, postIdInt, username, 1)
		fmt.Println(username + " upvoted" + " on post " + postId + " at " + now)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Upvote added"))
	case -1:
		if databaseAPI.HasDownvoted(database, username, postIdInt) {
			databaseAPI.RemoveVote(database, postIdInt, username)
			databaseAPI.DecreaseDownvotes(database, postIdInt)
			fmt.Println("Removed downvote from " + username + " on post " + postId + " at " + now)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Vote removed"))
			return
		}
		if databaseAPI.HasUpvoted(database, username, postIdInt) {
			databaseAPI.DecreaseUpvotes(database, postIdInt)
			databaseAPI.IncreaseDownvotes(database, postIdInt)
			databaseAPI.UpdateVote(database, postIdInt, username, -1)
			fmt.Println(username + " downvoted" + " on post " + postId + " at " + now)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Downvote added"))
			return
		}
		databaseAPI.IncreaseDownvotes(database, postIdInt)
		databaseAPI.AddVote(database, postIdInt, username, -1)
		fmt.Println(username + " downvoted" + " on post " + postId + " at " + now)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Downvote added"))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid vote"))
	}
}

// CommentsVoteApi handles voting on comments
func CommentsVoteApi(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !isLoggedIn(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	cookie, _ := r.Cookie("SESSION")
	username := databaseAPI.GetUser(database, cookie.Value)
	Id := r.FormValue("Id")
	IdInt, _ := strconv.Atoi(Id)
	vote := r.FormValue("vote")
	voteInt, _ := strconv.Atoi(vote)
	now := time.Now().Format("2006-01-02 15:04:05")

	switch voteInt {
	case 1:
		if databaseAPI.HasCommentUpvoted(database, username, IdInt) {
			databaseAPI.RemoveCommentVote(database, IdInt, username)
			databaseAPI.DecreaseCommentUpvotes(database, IdInt)
			fmt.Printf("Removed upvote from %s on post %d at %s\n", username, IdInt, now)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Vote removed"))
			return
		}
		if databaseAPI.HasCommentDownvoted(database, username, IdInt) {
			databaseAPI.DecreaseCommentDownvotes(database, IdInt)
			databaseAPI.IncreaseCommentUpvotes(database, IdInt)
			databaseAPI.UpdateCommentVote(database, IdInt, username, 1)
			fmt.Printf("%s changed vote from downvote to upvote on comment %d at %s\n", username, IdInt, now)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Upvote added"))
			return
		}
		databaseAPI.IncreaseCommentUpvotes(database, IdInt)
		databaseAPI.AddCommentVote(database, IdInt, username, 1)
		fmt.Printf("%s upvoted post %d at %s\n", username, IdInt, now)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Upvote added"))
	case -1:
		if databaseAPI.HasCommentDownvoted(database, username, IdInt) {
			databaseAPI.RemoveCommentVote(database, IdInt, username)
			databaseAPI.DecreaseCommentDownvotes(database, IdInt)
			fmt.Printf("Removed downvote from %s on post %d at %s\n", username, IdInt, now)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Vote removed"))
			return
		}
		if databaseAPI.HasCommentUpvoted(database, username, IdInt) {
			databaseAPI.DecreaseCommentUpvotes(database, IdInt)
			databaseAPI.IncreaseCommentDownvotes(database, IdInt)
			databaseAPI.UpdateCommentVote(database, IdInt, username, -1)
			fmt.Printf("%s changed vote from upvote to downvote on comment %d at %s\n", username, IdInt, now)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Downvote added"))
			return
		}
		databaseAPI.IncreaseCommentDownvotes(database, IdInt)
		databaseAPI.AddCommentVote(database, IdInt, username, -1)
		fmt.Printf("%s downvoted post %d at %s\n", username, IdInt, now)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Downvote added"))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid vote"))
	}
}

// DeletePostApi deletes a post by ID
func DeletePostApi(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Parse form values
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("ParseForm() error: %v", err), http.StatusInternalServerError)
		return
	}

	cookie, _ := r.Cookie("SESSION")
	username := databaseAPI.GetUser(database, cookie.Value)
	postId := r.FormValue("postId")
	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid postId: %v", err), http.StatusBadRequest)
		return
	}

	// Delete the post
	err = databaseAPI.DeletePost(database, postIdInt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting post: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Post %d deleted by %s\n", postIdInt, username)
	http.Redirect(w, r, "/filter?by=myposts", http.StatusFound)
}
