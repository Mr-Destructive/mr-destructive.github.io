package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// FrontMatter represents the front matter of a post
type FrontMatter struct {
	Title       string                 `json:"title" yaml:"title"`
	Description string                 `json:"description" yaml:"description"`
	Status      string                 `json:"status" yaml:"status"`
	Type        string                 `json:"type" yaml:"type"`
	Date        string                 `json:"date" yaml:"date"`
	Slug        string                 `json:"slug" yaml:"slug"`
	Tags        []string               `json:"tags" yaml:"tags"`
	ImageURL    string                 `json:"image_url" yaml:"image_url"`
	Extras      map[string]interface{} `json:",inline" yaml:",inline"`
}

// CleanFrontmatter cleans up inconsistent frontmatter in posts
func CleanFrontmatter(content []byte) ([]byte, error) {
	// Try to detect JSON front matter
	jsonSeparator := []byte("}\n\n")
	jsonIndex := indexOfBytes(content, jsonSeparator)

	if jsonIndex != -1 {
		frontmatterBytes := content[:jsonIndex+1] // Keep closing brace
		contentBytes := content[jsonIndex+2:]     // Skip the separator

		// Parse frontmatter
		var frontmatterObj FrontMatter
		tempMap := make(map[string]interface{})
		if err := json.Unmarshal(frontmatterBytes, &tempMap); err == nil {
			// Extract known fields into the struct
			if err := json.Unmarshal(frontmatterBytes, &frontmatterObj); err != nil {
				return nil, fmt.Errorf("error parsing JSON front matter: %v", err)
			}

			// Clean the frontmatter
			frontmatterObj = cleanFrontmatterFields(frontmatterObj)

			// Convert back to JSON
			cleanedFrontmatter, err := json.Marshal(frontmatterObj)
			if err != nil {
				return nil, fmt.Errorf("error marshaling cleaned frontmatter: %v", err)
			}

			// Combine cleaned frontmatter with content
			result := append(cleanedFrontmatter, []byte("\n\n")...)
			result = append(result, contentBytes...)
			return result, nil
		}
	}

	// Try to detect YAML front matter
	yamlSeparator := []byte("---\n\n")
	yamlIndex := indexOfBytes(content, yamlSeparator)

	if yamlIndex != -1 {
		frontmatterBytes := content[:yamlIndex]
		contentBytes := content[yamlIndex+len(yamlSeparator):]

		// Parse frontmatter
		var frontmatterObj FrontMatter
		tempMap := make(map[string]interface{})
		if err := yaml.Unmarshal(frontmatterBytes, &tempMap); err == nil {
			// Extract known fields into the struct
			if err := yaml.Unmarshal(frontmatterBytes, &frontmatterObj); err != nil {
				return nil, fmt.Errorf("error parsing YAML front matter: %v", err)
			}

			// Clean the frontmatter
			frontmatterObj = cleanFrontmatterFields(frontmatterObj)

			// Convert back to YAML
			cleanedFrontmatter, err := yaml.Marshal(frontmatterObj)
			if err != nil {
				return nil, fmt.Errorf("error marshaling cleaned frontmatter: %v", err)
			}

			// Combine cleaned frontmatter with content
			result := append(cleanedFrontmatter, []byte("---\n\n")...)
			result = append(result, contentBytes...)
			return result, nil
		}
	}

	return content, nil
}

// cleanFrontmatterFields cleans individual frontmatter fields
func cleanFrontmatterFields(fm FrontMatter) FrontMatter {
	// Ensure title is present
	if fm.Title == "" {
		fm.Title = "Untitled Post"
	}

	// Ensure status is valid
	if fm.Status != "published" && fm.Status != "draft" {
		fm.Status = "published"
	}

	// Ensure type is valid
	validTypes := map[string]bool{
		"posts":      true,
		"newsletter": true,
		"thoughts":   true,
		"projects":   true,
		"til":        true,
		"work":       true,
		"sqlog":      true,
	}
	if !validTypes[fm.Type] {
		fm.Type = "posts"
	}

	// Ensure date is in correct format
	if fm.Date != "" {
		// Try to parse the date
		_, err := time.Parse("2006-01-02", fm.Date)
		if err != nil {
			// If parsing fails, try other common formats
			possibleFormats := []string{
				"2006-1-2",
				"2006/01/02",
				"2006/1/2",
				"01/02/2006",
				"1/2/2006",
				"02/01/2006",
				"2/1/2006",
			}
			parsed := false
			for _, format := range possibleFormats {
				if date, err := time.Parse(format, fm.Date); err == nil {
					fm.Date = date.Format("2006-01-02")
					parsed = true
					break
				}
			}
			// If still not parsed, use current date
			if !parsed {
				fm.Date = time.Now().Format("2006-01-02")
			}
		}
	} else {
		// If no date, use current date
		fm.Date = time.Now().Format("2006-01-02")
	}

	// Ensure slug is present
	if fm.Slug == "" {
		fm.Slug = slugify(fm.Title)
	}

	// Ensure tags is a slice
	if fm.Tags == nil {
		fm.Tags = []string{}
	}

	return fm
}

// slugify creates a URL-friendly slug
func slugify(input string) string {
	// Convert to lowercase
	slug := strings.ToLower(input)
	
	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile("[^a-z0-9]+")
	slug = reg.ReplaceAllString(slug, "-")
	
	// Trim leading/trailing hyphens
	slug = strings.Trim(slug, "-")
	
	// If slug is empty, use a default
	if slug == "" {
		slug = "untitled"
	}
	
	return slug
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

// processFile cleans the frontmatter of a single file
func processFile(filePath string) error {
	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %v", filePath, err)
	}

	// Clean frontmatter
	cleanedContent, err := CleanFrontmatter(content)
	if err != nil {
		return fmt.Errorf("error cleaning frontmatter in file %s: %v", filePath, err)
	}

	// Write back to file only if content has changed
	if string(content) != string(cleanedContent) {
		err = os.WriteFile(filePath, cleanedContent, 0644)
		if err != nil {
			return fmt.Errorf("error writing file %s: %v", filePath, err)
		}
		fmt.Printf("Cleaned frontmatter in %s\n", filePath)
	}

	return nil
}

// processDirectory recursively processes all markdown files in a directory
func processDirectory(dirPath string) error {
	return filepath.WalkDir(dirPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".md" {
			err := processFile(path)
			if err != nil {
				fmt.Printf("Error processing file %s: %v\n", path, err)
			}
		}
		return nil
	})
}

func main() {
	// Process all posts directories
	dirs := []string{
		"posts/newsletter",
		"posts/posts",
		"posts/projects",
		"posts/thoughts",
		"posts/til",
		"posts/work",
		"posts/sqlog",
	}

	for _, dir := range dirs {
		fmt.Printf("Processing directory: %s\n", dir)
		err := processDirectory(dir)
		if err != nil {
			fmt.Printf("Error processing directory %s: %v\n", dir, err)
		}
	}

	fmt.Println("Frontmatter cleanup completed!")
}