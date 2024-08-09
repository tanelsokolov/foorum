# Forum

The Literary Lions forum is a straightforward discussion platform with an intuitive interface. Users can create topics, categorize them, reply to posts, search posts, and upvote or downvote topics.

## Backend

Our backend is built with Golang and interacts with an SQLite database to manage data and handle user authentication. The database stores:

- User information
- Topics
- Replies
- Votes

## Frontend

The frontend is developed using HTML, CSS, and JavaScript, connecting to the backend to interact with the database.

## Authentication

Upon login, users receive a UUID token in a session cookie for authentication. This token is:

- Stored in the database
- Used for user verification

When a user logs out, the token is removed from the database. During registration, the username and a bcrypt-hashed password are stored in the database.

## Communication

- **Logged-in users** can:
  - Create posts
  - Add comments
  - Reply to posts
  - View topics and comments
  - Delete topics and comments

- **Non-logged-in users** can:
  - View topics and comments

## Like and Dislike

- **Logged-in users** can:
  - Upvote/downvote posts/comments

- **Non-logged-in users** can:
  - Only see counts of votes

## Filter Posts

- **Logged-in users** can:
  - Filter posts by categories
  - Filter posts by created posts
  - Filter posts by liked posts

- **Non-logged-in users** can:
  - Only filter posts by categories

## Docker

To use Docker to run the application:

1. Create a `Dockerfile` in the root directory of the repository.
2. Build the Docker image with the following command:

   ```bash
   docker build -t lions .
   
3. Run the application with the following command:
   ```bash
   docker run -p 8000:8000 -it lions

The server will be running at http://localhost:8000/.