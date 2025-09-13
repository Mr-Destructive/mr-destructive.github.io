package main

import (
	"net/http"

	"github.com/mr-destructive/mr-destructive.github.io/plugins"
)

func main() {

	http.HandleFunc("/editor", plugins.PostHandler)
	http.ListenAndServe(":8081", nil)
}
