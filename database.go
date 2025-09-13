package main

import (
	\"database/sql\"
	\"encoding/json\"
	\"fmt\"
	\"log\"
	\"os\"
	\"path/filepath\"
	\"sort\"
	\"strings\"
	\"time\"

	_ \"github.com/mattn/go-sqlite3\"
	models \"github.com/mr-destructive/mr-destructive.github.io/models\"
)

// DB represents the database connection
type DB struct {
	*sql.DB
}

// InitDB initializes the SQLite database
func InitDB(dataSourceName string) (*DB, error) {
	database, err := sql.Open(\"sqlite3\", dataSourceName)
	if err != nil {
		return nil, err
	}
	if err := database.Ping(); err != nil {
		return nil, err
	}
	return &DB{database}, nil
}

// CreateTables creates the necessary tables
func (db *DB) CreateTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS authors (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		password TEXT NOT NULL,
		is_admin BOOLEAN DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		slug TEXT NOT NULL UNIQUE,
		body TEXT NOT NULL,
		metadata TEXT NOT NULL,
		deleted BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		author_id INTEGER REFERENCES authors(id) NOT NULL
	);
	`
	_, err := db.Exec(query)
	return err
}

// InsertAuthor inserts a new author
func (db *DB) InsertAuthor(username, name, password string) (int64, error) {
	result, err := db.Exec(\"INSERT INTO authors (username, name, password) VALUES (?, ?, ?)\", username, name, password)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// InsertPost inserts a new post
func (db *DB) InsertPost(title, slug, body, metadata string, authorID int64) (int64, error) {
	result, err := db.Exec(\"INSERT INTO posts (title, slug, body, metadata, author_id) VALUES (?, ?, ?, ?, ?)\", title, slug, body, metadata, authorID)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetPostByID retrieves a post by ID
func (db *DB) GetPostByID(id int64) (*models.DBPost, error) {
	row := db.QueryRow(\"SELECT id, title, slug, body, metadata, deleted, created_at, updated_at, author_id FROM posts WHERE id = ?\", id)
	var post models.DBPost
	err := row.Scan(&post.ID, &post.Title, &post.Slug, &post.Body, &post.Metadata, &post.Deleted, &post.CreatedAt, &post.UpdatedAt, &post.AuthorID)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetPostBySlug retrieves a post by slug
func (db *DB) GetPostBySlug(slug string) (*models.DBPost, error) {
	row := db.QueryRow(\"SELECT id, title, slug, body, metadata, deleted, created_at, updated_at, author_id FROM posts WHERE slug = ?\", slug)
	var post models.DBPost
	err := row.Scan(&post.ID, &post.Title, &post.Slug, &post.Body, &post.Metadata, &post.Deleted, &post.CreatedAt, &post.UpdatedAt, &post.AuthorID)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// UpdatePost updates a post
func (db *DB) UpdatePost(id int64, title, slug, body, metadata string) error {
	_, err := db.Exec(\"UPDATE posts SET title = ?, slug = ?, body = ?, metadata = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?\", title, slug, body, metadata, id)
	return err
}

// DeletePost marks a post as deleted
func (db *DB) DeletePost(id int64) error {
	_, err := db.Exec(\"UPDATE posts SET deleted = 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?\", id)
	return err
}

// GetAllPosts retrieves all non-deleted posts
func (db *DB) GetAllPosts() ([]models.DBPost, error) {
	rows, err := db.Query(\"SELECT id, title, slug, body, metadata, deleted, created_at, updated_at, author_id FROM posts WHERE deleted = 0 ORDER BY created_at ASC\")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []models.DBPost
	for rows.Next() {
		var post models.DBPost
		err := rows.Scan(&post.ID, &post.Title, &post.Slug, &post.Body, &post.Metadata, &post.Deleted, &post.CreatedAt, &post.UpdatedAt, &post.AuthorID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

// SyncPostsToDB syncs all markdown posts to the database
func (db *DB) SyncPostsToDB(posts []models.Post, authorID int64) error {
	// First, get all existing posts from DB
	existingPosts, err := db.GetAllPosts()
	if err != nil {
		return err
	}

	// Create a map of existing slugs for quick lookup
	existingSlugs := make(map[string]bool)
	for _, post := range existingPosts {
		existingSlugs[post.Slug] = true
	}

	// Sort posts by date to ensure consistent ID assignment
	sort.Slice(posts, func(i, j int) bool {
		date1, err1 := time.Parse(\"2006-01-02\", posts[i].Frontmatter.Date)
		date2, err2 := time.Parse(\"2006-01-02\", posts[j].Frontmatter.Date)
		if err1 != nil || err2 != nil {
			return false
		}
		return date1.Before(date2)
	})

	// Insert new posts
	for _, post := range posts {
		slug := post.Frontmatter.Slug
		// Skip if already exists
		if existingSlugs[slug] {
			continue
		}

		// Convert metadata to JSON
		metadataBytes, err := json.Marshal(post.Frontmatter)
		if err != nil {
			log.Printf(\"Error marshaling metadata for post %s: %v\", post.Frontmatter.Title, err)
			continue
		}

		_, err = db.InsertPost(
			post.Frontmatter.Title,
			slug,
			post.Markdown,
			string(metadataBytes),
			authorID,
		)
		if err != nil {
			log.Printf(\"Error inserting post %s: %v\", post.Frontmatter.Title, err)
			continue
		}
		log.Printf(\"Inserted post: %s\", post.Frontmatter.Title)
	}

	return nil
}

// CleanFrontmatter cleans up inconsistent frontmatter
func CleanFrontmatter(posts []models.Post) []models.Post {
	var cleanedPosts []models.Post
	for _, post := range posts {
		// Ensure required fields are present
		if post.Frontmatter.Title == \"\" {
			// Try to extract title from content if possible
			lines := strings.Split(post.Markdown, \"\\n\")
			if len(lines) > 0 {
				// Remove markdown headers (#) if present
				title := strings.TrimPrefix(lines[0], \"# \")
				title = strings.TrimSpace(title)
				if title != \"\" {
					post.Frontmatter.Title = title
				}
			}
		}

		// Ensure date is in correct format
		if post.Frontmatter.Date != \"\" {
			// Try to parse the date
			_, err := time.Parse(\"2006-01-02\", post.Frontmatter.Date)
			if err != nil {
				// If parsing fails, try other common formats
				possibleFormats := []string{
					\"2006-1-2\",
					\"2006/01/02\",
					\"2006/1/2\",
					\"01/02/2006\",
					\"1/2/2006\",
					\"02/01/2006\",
					\"2/1/2006\",
				}
				parsed := false
				for _, format := range possibleFormats {
					if date, err := time.Parse(format, post.Frontmatter.Date); err == nil {
						post.Frontmatter.Date = date.Format(\"2006-01-02\")
						parsed = true
						break
					}
				}
				// If still not parsed, use current date
				if !parsed {
					post.Frontmatter.Date = time.Now().Format(\"2006-01-02\")
				}
			}
		} else {
			// If no date, use current date
			post.Frontmatter.Date = time.Now().Format(\"2006-01-02\")
		}

		// Ensure type is set
		if post.Frontmatter.Type == \"\" {
			post.Frontmatter.Type = \"posts\"
		}

		// Ensure slug is set
		if post.Frontmatter.Slug == \"\" {
			post.Frontmatter.Slug = Slugify(post.Frontmatter.Title)
		}

		// Ensure status is set
		if post.Frontmatter.Status == \"\" {
			post.Frontmatter.Status = \"published\"
		}

		cleanedPosts = append(cleanedPosts, post)
	}
	return cleanedPosts
}

// Slugify creates a URL-friendly slug
func Slugify(input string) string {
	// Replace spaces and special characters with hyphens
	slug := strings.ToLower(input)
	slug = strings.ReplaceAll(slug, \" \", \"-\")
	slug = strings.ReplaceAll(slug, \"--\", \"-\")
	// Remove any non-alphanumeric characters except hyphens
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)
	// Trim leading/trailing hyphens
	slug = strings.Trim(slug, \"-\")
	return slug
}