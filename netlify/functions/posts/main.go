package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/mattn/go-sqlite3"
)

// DB represents the database connection
type DB struct {
	*sql.DB
}

// InitDB initializes the SQLite database
func InitDB(dataSourceName string) (*DB, error) {
	database, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := database.Ping(); err != nil {
		return nil, err
	}
	return &DB{database}, nil
}

// Post represents a blog post
type Post struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Slug      string `json:"slug"`
	Body      string `json:"body"`
	Metadata  string `json:"metadata"`
	Deleted   bool   `json:"deleted"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	AuthorID  int    `json:"author_id"`
}

// GetPostByID retrieves a post by ID
func (db *DB) GetPostByID(id int) (*Post, error) {
	row := db.QueryRow("SELECT id, title, slug, body, metadata, deleted, created_at, updated_at, author_id FROM posts WHERE id = ? AND deleted = 0", id)
	var post Post
	err := row.Scan(&post.ID, &post.Title, &post.Slug, &post.Body, &post.Metadata, &post.Deleted, &post.CreatedAt, &post.UpdatedAt, &post.AuthorID)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetAllPosts retrieves all non-deleted posts
func (db *DB) GetAllPosts() ([]Post, error) {
	rows, err := db.Query("SELECT id, title, slug, body, metadata, deleted, created_at, updated_at, author_id FROM posts WHERE deleted = 0 ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Body, &post.Metadata, &post.Deleted, &post.CreatedAt, &post.UpdatedAt, &post.AuthorID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// InsertPost inserts a new post
func (db *DB) InsertPost(title, slug, body, metadata string, authorID int) (int, error) {
	result, err := db.Exec("INSERT INTO posts (title, slug, body, metadata, author_id) VALUES (?, ?, ?, ?, ?)", title, slug, body, metadata, authorID)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// UpdatePost updates a post
func (db *DB) UpdatePost(id int, title, slug, body, metadata string) error {
	_, err := db.Exec("UPDATE posts SET title = ?, slug = ?, body = ?, metadata = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", title, slug, body, metadata, id)
	return err
}

// DeletePost marks a post as deleted
func (db *DB) DeletePost(id int) error {
	_, err := db.Exec("UPDATE posts SET deleted = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?", id)
	return err
}

// Handler is the main Lambda handler
func Handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Initialize database
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./data/blog.db"
	}
	
	db, err := InitDB(dbPath)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("{\"error\": \"Database connection failed: %s\"}", err.Error()),
		}, nil
	}
	defer db.Close()

	// Handle CORS preflight requests
	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
			Body: "",
		}, nil
	}

	// Add CORS headers to all responses
	headers := map[string]string{
		"Access-Control-Allow-Origin":  "*",
		"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
		"Access-Control-Allow-Headers": "Content-Type, Authorization",
		"Content-Type":                 "application/json",
	}

	// Route based on path and method
	path := req.Path
	method := req.HTTPMethod

	// Handle /posts endpoint
	if strings.HasPrefix(path, "/posts") {
		parts := strings.Split(strings.Trim(path, "/"), "/")
		
		// /posts - GET: list all posts, POST: create new post
		if len(parts) == 1 {
			switch method {
			case "GET":
				return listPosts(db, headers)
			case "POST":
				return createPost(db, req, headers)
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusMethodNotAllowed,
					Headers:    headers,
					Body:       "{\"error\": \"Method not allowed\"}",
				}, nil
			}
		}
		
		// /posts/{id} - GET: get post, PUT: update post, DELETE: delete post
		if len(parts) == 2 {
			postID := parts[1]
			switch method {
			case "GET":
				return getPost(db, postID, headers)
			case "PUT":
				return updatePost(db, postID, req, headers)
			case "DELETE":
				return deletePost(db, postID, headers)
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: http.StatusMethodNotAllowed,
					Headers:    headers,
					Body:       "{\"error\": \"Method not allowed\"}",
				}, nil
			}
		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusNotFound,
		Headers:    headers,
		Body:       "{\"error\": \"Endpoint not found\"}",
	}, nil
}

func main() {
	lambda.Start(Handler)
}

// Helper functions for each operation
func listPosts(db *DB, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	posts, err := db.GetAllPosts()
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf("{\"error\": \"Failed to fetch posts: %s\"}", err.Error()),
		}, nil
	}

	jsonPosts, err := json.Marshal(posts)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf("{\"error\": \"Failed to marshal posts: %s\"}", err.Error()),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(jsonPosts),
	}, nil
}

func getPost(db *DB, postID string, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	// Parse post ID
	var id int
	_, err := fmt.Sscanf(postID, "%d", &id)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
			Body:       "{\"error\": \"Invalid post ID\"}",
		}, nil
	}

	post, err := db.GetPostByID(id)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Headers:    headers,
			Body:       fmt.Sprintf("{\"error\": \"Post not found: %s\"}", err.Error()),
		}, nil
	}

	jsonPost, err := json.Marshal(post)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf("{\"error\": \"Failed to marshal post: %s\"}", err.Error()),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(jsonPost),
	}, nil
}

func createPost(db *DB, req events.APIGatewayProxyRequest, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	// Parse request body
	var postData struct {
		Title    string `json:"title"`
		Slug     string `json:"slug"`
		Body     string `json:"body"`
		Metadata string `json:"metadata"`
		AuthorID int    `json:"author_id"`
	}
	
	err := json.Unmarshal([]byte(req.Body), &postData)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
			Body:       fmt.Sprintf("{\"error\": \"Invalid request body: %s\"}", err.Error()),
		}, nil
	}

	// Insert post
	postID, err := db.InsertPost(postData.Title, postData.Slug, postData.Body, postData.Metadata, postData.AuthorID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf("{\"error\": \"Failed to create post: %s\"}", err.Error()),
		}, nil
	}

	// Return created post
	post, err := db.GetPostByID(postID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf("{\"error\": \"Failed to fetch created post: %s\"}", err.Error()),
		}, nil
	}

	jsonPost, err := json.Marshal(post)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf("{\"error\": \"Failed to marshal post: %s\"}", err.Error()),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusCreated,
		Headers:    headers,
		Body:       string(jsonPost),
	}, nil
}

func updatePost(db *DB, postID string, req events.APIGatewayProxyRequest, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	// Parse post ID
	var id int
	_, err := fmt.Sscanf(postID, "%d", &id)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
			Body:       "{\"error\": \"Invalid post ID\"}",
		}, nil
	}

	// Parse request body
	var postData struct {
		Title    string `json:"title"`
		Slug     string `json:"slug"`
		Body     string `json:"body"`
		Metadata string `json:"metadata"`
	}
	
	err = json.Unmarshal([]byte(req.Body), &postData)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
			Body:       fmt.Sprintf("{\"error\": \"Invalid request body: %s\"}", err.Error()),
		}, nil
	}

	// Update post
	err = db.UpdatePost(id, postData.Title, postData.Slug, postData.Body, postData.Metadata)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf("{\"error\": \"Failed to update post: %s\"}", err.Error()),
		}, nil
	}

	// Return updated post
	post, err := db.GetPostByID(id)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf("{\"error\": \"Failed to fetch updated post: %s\"}", err.Error()),
		}, nil
	}

	jsonPost, err := json.Marshal(post)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf("{\"error\": \"Failed to marshal post: %s\"}", err.Error()),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       string(jsonPost),
	}, nil
}

func deletePost(db *DB, postID string, headers map[string]string) (events.APIGatewayProxyResponse, error) {
	// Parse post ID
	var id int
	_, err := fmt.Sscanf(postID, "%d", &id)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
			Body:       "{\"error\": \"Invalid post ID\"}",
		}, nil
	}

	// Delete post
	err = db.DeletePost(id)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
			Body:       fmt.Sprintf("{\"error\": \"Failed to delete post: %s\"}", err.Error()),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers,
		Body:       "{\"message\": \"Post deleted successfully\"}",
	}, nil
}