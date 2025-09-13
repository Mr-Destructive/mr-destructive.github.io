package main

import (
	\"database/sql\"
	"encoding/json\"
	"fmt\"
	"log\"
	"os\"
	"path/filepath\"

	_ \"github.com/mattn/go-sqlite3\"
	models \"github.com/mr-destructive/mr-destructive.github.io/models\"
	\"gopkg.in/yaml.v3\"
)

func main() {
	// Initialize SQLite database
	db, err := InitDB(\"./data/blog.db\")
	if err != nil {
		log.Fatal(\"Failed to initialize database:\", err)
	}
	defer db.Close()

	// Create tables
	err = db.CreateTables()
	if err != nil {
		log.Fatal(\"Failed to create tables:\", err)
	}

	// Create default author if not exists
	authorID, err := db.InsertAuthor(\"admin\", \"Admin User\", \"admin123\")
	if err != nil {
		// If author already exists, get the existing author ID
		var existingID int64
		err = db.QueryRow(\"SELECT id FROM authors WHERE username = ?\", \"admin\").Scan(&existingID)
		if err != nil {
			log.Fatal(\"Failed to get existing author:\", err)
		}
		authorID = existingID
	}

	// Read all posts
	configBytes, err := os.ReadFile(models.SSG_CONFIG_FILE_NAME)
	if err != nil {
		log.Fatal(\"Failed to read config:\", err)
	}

	var config models.SSG_CONFIG
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatal(\"Failed to parse config:\", err)
	}

	// Walk through posts directory
	postFiles, err := WalkAndListFiles(config.Blog.PostsDir)
	if err != nil {
		log.Fatal(\"Failed to list post files:\", err)
	}

	posts, err := ReadPosts(postFiles)
	if err != nil {
		log.Fatal(\"Failed to read posts:\", err)
	}

	// Clean frontmatter
	cleanedPosts := CleanFrontmatter(posts)

	// Sync posts to database
	err = db.SyncPostsToDB(cleanedPosts, authorID)
	if err != nil {
		log.Fatal(\"Failed to sync posts to database:\", err)
	}

	fmt.Println(\"Database initialized and posts synced successfully!\")
}

// WalkAndListFiles walks through a directory and lists all files
func WalkAndListFiles(dirPath string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == \".md\" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// ReadFiles reads the content of multiple files
func ReadFiles(files []string) ([][]byte, error) {
	var filesBytes = [][]byte{}
	for _, file := range files {
		fileBytes, err := os.ReadFile(file)
		if err != nil {
			return filesBytes, err
		}
		filesBytes = append(filesBytes, fileBytes)
	}
	return filesBytes, nil
}

// ReadPosts reads and parses markdown posts
func ReadPosts(files []string) ([]models.Post, error) {
	var posts []models.Post

	// Read file contents
	filesBytes, err := ReadFiles(files)
	if err != nil {
		return nil, err
	}

	// Iterate through files
	for _, fileBytes := range filesBytes {
		var success bool
		var frontmatterObj models.FrontMatter
		var contentBytes []byte
		var requiredFields []string = []string{\"title\", \"description\", \"status\", \"type\", \"date\", \"slug\", \"tags\", \"image_url\"}

		// Attempt to detect JSON front matter
		jsonSeparator := []byte(\"}\\n\\n\")
		jsonIndex := indexOfBytes(fileBytes, jsonSeparator)

		if jsonIndex != -1 {
			frontmatterBytes := fileBytes[:jsonIndex+1] // Keep closing brace
			contentBytes = fileBytes[jsonIndex+2:]      // Skip the separator

			// Unmarshal into a temporary map to capture extra fields
			tempMap := make(map[string]interface{})
			if err := json.Unmarshal(frontmatterBytes, &tempMap); err == nil {
				success = true

				// Extract known fields into the struct
				if err := json.Unmarshal(frontmatterBytes, &frontmatterObj); err != nil {
					log.Printf(\"Error parsing JSON front matter: %v\", err)
					continue
				}

				// Remove known keys and store the rest in Extras
				for _, key := range requiredFields {
					delete(tempMap, key)
				}
				frontmatterObj.Extras = tempMap
			}
		}

		// Attempt to detect YAML front matter
		if !success {
			yamlSeparator := []byte(\"---\\n\\n\")
			yamlIndex := indexOfBytes(fileBytes, yamlSeparator)

			if yamlIndex != -1 {
				frontmatterBytes := fileBytes[:yamlIndex]
				contentBytes = fileBytes[yamlIndex+len(yamlSeparator):]

				// Unmarshal into a temporary map to capture extra fields
				tempMap := make(map[string]interface{})
				if err := yaml.Unmarshal(frontmatterBytes, &tempMap); err == nil {

					// Extract known fields into the struct
					if err := yaml.Unmarshal(frontmatterBytes, &frontmatterObj); err != nil {
						log.Printf(\"Error parsing YAML front matter: %v\", err)
						continue
					}

					// Remove known keys and store the rest in Extras
					for _, key := range requiredFields {
						delete(tempMap, key)
					}
					frontmatterObj.Extras = tempMap
				} else {
					log.Printf(\"Error parsing YAML front matter: %v\", err)
					continue
				}
			} else {
				log.Printf(\"No valid front matter found in file\")
				continue
			}
		}

		// Append post
		posts = append(posts, models.Post{
			Frontmatter: frontmatterObj,
			Markdown:    string(contentBytes),
		})
	}

	return posts, nil
}

// indexOfBytes finds the index of a byte sequence in a byte slice
func indexOfBytes(slice, sub []byte) int {
	for i := 0; i <= len(slice)-len(sub); i++ {
		if string(slice[i:i+len(sub)]) == string(sub) {
			return i
		}
	}
	return -1
}