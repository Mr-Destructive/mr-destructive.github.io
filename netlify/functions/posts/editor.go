package main

import (
	"bytes"
	"html/template"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// TemplateData holds data for the template
type TemplateData struct {
	Themes struct {
		Default struct {
			Bg            string `json:"bg"`
			Text          string `json:"text"`
			SecondaryText string `json:"secondary-text"`
			Link          struct {
				Normal string `json:"normal"`
				Hover  string `json:"hover"`
				Active string `json:"active"`
			} `json:"link"`
			Quotes     string `json:"quotes"`
			CodeBlocks struct {
				Bg     string `json:"bg"`
				Border string `json:"border"`
			} `json:"codeblocks"`
			Code struct {
				Text     string `json:"text"`
				Comment  string `json:"comment"`
				Keyword  string `json:"keyword"`
				String   string `json:"string"`
				Number   string `json:"number"`
				Variable string `json:"variable"`
				Function string `json:"function"`
			} `json:"code"`
		} `json:"default"`
		Secondary struct {
			Bg            string `json:"bg"`
			Text          string `json:"text"`
			SecondaryText string `json:"secondary-text"`
			Link          struct {
				Normal string `json:"normal"`
				Hover  string `json:"hover"`
				Active string `json:"active"`
			} `json:"link"`
			Quotes     string `json:"quotes"`
			CodeBlocks struct {
				Bg     string `json:"bg"`
				Border string `json:"border"`
			} `json:"codeblocks"`
			Code struct {
				Text     string `json:"text"`
				Comment  string `json:"comment"`
				Keyword  string `json:"keyword"`
				String   string `json:"string"`
				Number   string `json:"number"`
				Variable string `json:"variable"`
				Function string `json:"function"`
			} `json:"code"`
		} `json:"secondary"`
	} `json:"themes"`
	Config struct {
		Blog struct {
			PrefixURL string `json:"prefix_url"`
		} `json:"blog"`
	} `json:"config"`
}

// EditorHandler serves the editor template
func EditorHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
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

	// Read the editor template file
	templateBytes, err := os.ReadFile("editor.html")
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
				"Content-Type":                 "text/html",
			},
			Body: "<html><body><h1>Error loading editor</h1><p>Could not load editor template.</p></body></html>",
		}, nil
	}

	// Parse the template
	tmpl, err := template.New("editor").Parse(string(templateBytes))
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
				"Content-Type":                 "text/html",
			},
			Body: "<html><body><h1>Error parsing template</h1><p>Could not parse editor template.</p></body></html>",
		}, nil
	}

	// Prepare template data
	data := TemplateData{}
	
	// Set default theme values
	data.Themes.Default.Bg = "#ffffff"
	data.Themes.Default.Text = "#333333"
	data.Themes.Default.SecondaryText = "#00ffff"
	data.Themes.Default.Link.Normal = "#007bff"
	data.Themes.Default.Link.Hover = "#0056b3"
	data.Themes.Default.Link.Active = "#003a75"
	data.Themes.Default.Quotes = "#999999"
	data.Themes.Default.CodeBlocks.Bg = "#dddddd"
	data.Themes.Default.CodeBlocks.Border = "#ced4da"
	data.Themes.Default.Code.Text = "#444444"
	data.Themes.Default.Code.Comment = "#808080"
	data.Themes.Default.Code.Keyword = "#008000"
	data.Themes.Default.Code.String = "#000080"
	data.Themes.Default.Code.Number = "#000000"
	data.Themes.Default.Code.Variable = "#000000"
	data.Themes.Default.Code.Function = "#000000"
	
	// Set secondary theme values
	data.Themes.Secondary.Bg = "#121212"
	data.Themes.Secondary.Text = "#ffffff"
	data.Themes.Secondary.SecondaryText = "#00ffff"
	data.Themes.Secondary.Link.Normal = "#ff6600"
	data.Themes.Secondary.Link.Hover = "#4682b4"
	data.Themes.Secondary.Link.Active = "#00008b"
	data.Themes.Secondary.Quotes = "#a9a9a9"
	data.Themes.Secondary.CodeBlocks.Bg = "#333333"
	data.Themes.Secondary.CodeBlocks.Border = "#444444"
	data.Themes.Secondary.Code.Text = "#ffffff"
	data.Themes.Secondary.Code.Comment = "#b0b0b0"
	data.Themes.Secondary.Code.Keyword = "#32cd32"
	data.Themes.Secondary.Code.String = "#ff6347"
	data.Themes.Secondary.Code.Number = "#d3d3d3"
	data.Themes.Secondary.Code.Variable = "#b0e0e6"
	data.Themes.Secondary.Code.Function = "#ff4500"
	
	// Set config values
	data.Config.Blog.PrefixURL = ""

	// Execute the template with data
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
				"Content-Type":                 "text/html",
			},
			Body: "<html><body><h1>Error executing template</h1><p>Could not execute editor template.</p></body></html>",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type, Authorization",
			"Content-Type":                 "text/html",
		},
		Body: buf.String(),
	}, nil
}

func main() {
	lambda.Start(EditorHandler)
}