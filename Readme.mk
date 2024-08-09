# Forum

Literary Lions forum is a straightforward discussion platform with an intuitive interface. 
Users can create topics, categorize them, reply to posts, search from posts and upvote or downvote topics.

## Backend

Our backend is built with Golang and interacts with an SQLite database to manage data and handle user authentication.
The database stores user information, topics, replies, and votes.

## Frontend

The frontend is developed using HTML, CSS, and JavaScript, connecting to the backend to interact with the database.

## Authentication

Upon login, users receive a UUID token in a session cookie for authentication.
This token is stored in the database and used for user verification. When a user logs out, the token is removed from the database.
During registration, the username and a bcrypt-hashed password are stored in the database.

## Communication

Users that are logged in, can create posts, add comments, reply to posts, view topics and comments, delete topics and comments.
Users that are not logged in, can only view topics and comments

## Like and dislike

Users that are logged in can upvote/downvote posts/comments.
Users that are not logged in, can only see counts of votes.

## Filter posts

Users that are logged in can filter posts by categories, created posts and liked posts.
Users that are not logged in can only filter posts by categories.

## Docker

To use Docker to run the application, create a Dockerfile in the root directory of the repository.

To create the Docker image, we run the following command:
```bash
docker build -t lions .
```
To run the application, we use the following command:
```bash
docker run -p 8000:8000 -it lions
```
The server will be running at http://localhost:8000/

## When not using Docker to run the application

```bash
go run main.go
```
The server will be running at http://localhost:8000/