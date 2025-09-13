package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yuin/goldmark"
	"bytes"
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

// ViewHandler handles viewing a post by ID
func ViewHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Handle CORS preflight requests
	if req.HTTPMethod == "OPTIONS" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
			Body: "",
		}, nil
	}

	// Get post ID from path parameter
	pathParts := strings.Split(req.Path, "/")
	if len(pathParts) < 3 {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
				"Content-Type":                 "text/html",
			},
			Body: "<html><body><h1>Bad Request</h1><p>Missing post ID.</p></body></html>",
		}, nil
	}

	postIDStr := pathParts[2]
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
				"Content-Type":                 "text/html",
			},
			Body: "<html><body><h1>Bad Request</h1><p>Invalid post ID.</p></body></html>",
		}, nil
	}

	// Initialize database
	dbPath := os.Getenv("DATABASE_PATH")
	if dbPath == "" {
		dbPath = "./data/blog.db"
	}
	
	db, err := InitDB(dbPath)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
				"Content-Type":                 "text/html",
			},
			Body: fmt.Sprintf("<html><body><h1>Database Error</h1><p>Failed to connect to database: %s</p></body></html>", err.Error()),
		}, nil
	}
	defer db.Close()

	// Get post
	post, err := db.GetPostByID(postID)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
				"Content-Type":                 "text/html",
			},
			Body: "<html><body><h1>Not Found</h1><p>Post not found.</p></body></html>",
		}, nil
	}

	// Parse metadata
	var metadata map[string]interface{}
	err = json.Unmarshal([]byte(post.Metadata), &metadata)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
				"Content-Type":                 "text/html",
			},
			Body: fmt.Sprintf("<html><body><h1>Metadata Error</h1><p>Failed to parse metadata: %s</p></body></html>", err.Error()),
		}, nil
	}

	// Convert markdown to HTML
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(post.Body), &buf); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
				"Content-Type":                 "text/html",
			},
			Body: fmt.Sprintf("<html><body><h1>Markdown Error</h1><p>Failed to convert markdown: %s</p></body></html>", err.Error()),
		}, nil
	}
	contentHTML := buf.String()

	// Create HTML response
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
            line-height: 1.6;
            max-width: 800px;
            margin: 0 auto;
            padding: 2rem;
            color: #333;
            background-color: #fff;
        }
        h1, h2, h3 {
            color: #2c3e50;
        }
        code {
            background-color: #f4f4f4;
            padding: 0.2em 0.4em;
            border-radius: 3px;
            font-family: 'SF Mono', Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace;
        }
        pre {
            background-color: #f4f4f4;
            padding: 1rem;
            border-radius: 5px;
            overflow-x: auto;
        }
        blockquote {
            border-left: 4px solid #3498db;
            padding-left: 1rem;
            margin-left: 0;
            color: #7f8c8d;
        }
        a {
            color: #3498db;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
        .post-meta {
            color: #7f8c8d;
            font-size: 0.9rem;
            margin-bottom: 2rem;
        }
        .back-link {
            display: inline-block;
            margin-bottom: 1rem;
            color: #3498db;
        }
    </style>
</head>
<body>
    <a href="/" class="back-link">&larr; Back to Home</a>
    <h1>%s</h1>
    <div class="post-meta">
        <span>Published on %s</span>
        %s
    </div>
    <div class="post-content">
        %s
    </div>
</body>
</html>`,
		post.Title,
		post.Title,
		metadata["date"],
		func() string {
			if tags, ok := metadata["tags"].([]interface{}); ok && len(tags) > 0 {
				tagStrs := make([]string, len(tags))
				for i, tag := range tags {
					tagStrs[i] = fmt.Sprintf("%v", tag)
				}
				return fmt.Sprintf(" | Tags: %s", strings.Join(tagStrs, ", "))
			}
			return ""
		}(),
		contentHTML)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type, Authorization",
			"Content-Type":                 "text/html",
		},
		Body: html,
	}, nil
}

func main() {
	lambda.Start(ViewHandler)
}