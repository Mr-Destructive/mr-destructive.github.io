package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Metadata struct {
	Title  string   `json:"title"`
	Slug   string   `json:"slug"`
	Type   string   `json:"type"`
	Status string   `json:"status"`
	Author string   `json:"author"`
	Date   string   `json:"date"`
	Tags   []string `json:"tags"`
}

func main() {
	// Initialize SQLite database
	db, err := sql.Open("sqlite3", "./data/blog.db")
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Create sample author
	_, err = db.Exec("INSERT OR IGNORE INTO authors (username, name, password) VALUES (?, ?, ?)", "admin", "Admin User", "admin123")
	if err != nil {
		log.Fatal("Failed to create author:", err)
	}

	// Get author ID
	var authorID int
	err = db.QueryRow("SELECT id FROM authors WHERE username = ?", "admin").Scan(&authorID)
	if err != nil {
		log.Fatal("Failed to get author ID:", err)
	}

	// Create sample posts
	samplePosts := []struct {
		Title    string
		Slug     string
		Body     string
		Metadata Metadata
	}{
		{
			Title: "Welcome to My Blog",
			Slug:  "welcome-to-my-blog",
			Body:  "# Welcome\n\nThis is my first post on this blog. I'm excited to share my thoughts and experiences with you.\n\n## What to Expect\n\n- Technical tutorials\n- Personal reflections\n- Book reviews\n- Project showcases",
			Metadata: Metadata{
				Title:  "Welcome to My Blog",
				Slug:   "welcome-to-my-blog",
				Type:   "posts",
				Status: "published",
				Author: "Admin User",
				Date:   "2024-01-01",
				Tags:   []string{"welcome", "introduction"},
			},
		},
		{
			Title: "Building a Static Site Generator",
			Slug:  "building-static-site-generator",
			Body:  "# Building a Static Site Generator\n\nToday I'll walk you through how I built my own static site generator using Go.\n\n## Why Build Your Own?\n\n1. Learning experience\n2. Complete control over features\n3. Performance optimization\n\n```go\nfunc main() {\n    fmt.Println(\"Hello, SSG!\")\n}\n```\n\n## Key Components\n\n- File parsing\n- Template rendering\n- Asset handling\n- Deployment scripts",
			Metadata: Metadata{
				Title:  "Building a Static Site Generator",
				Slug:   "building-static-site-generator",
				Type:   "posts",
				Status: "published",
				Author: "Admin User",
				Date:   "2024-01-15",
				Tags:   []string{"go", "ssg", "programming"},
			},
		},
	}

	for _, post := range samplePosts {
		// Convert metadata to JSON
		metadataJSON, err := json.Marshal(post.Metadata)
		if err != nil {
			log.Printf("Failed to marshal metadata for post %s: %v", post.Title, err)
			continue
		}

		// Insert post
		_, err = db.Exec("INSERT OR IGNORE INTO posts (title, slug, body, metadata, author_id) VALUES (?, ?, ?, ?, ?)",
			post.Title, post.Slug, post.Body, string(metadataJSON), authorID)
		if err != nil {
			log.Printf("Failed to insert post %s: %v", post.Title, err)
			continue
		}

		fmt.Printf("Inserted post: %s\n", post.Title)
	}

	fmt.Println("Database populated with sample data!")
}