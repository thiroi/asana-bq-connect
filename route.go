package src

import (
"net/http"
"src/main"
)

func init() {
	// net/http
	http.HandleFunc("/01", main.Main)
}
