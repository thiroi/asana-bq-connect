package main

import (
	"net/http"
)

func HistoryConnector(w http.ResponseWriter, r *http.Request) {
	connect(w, r, "Sprint-180820", true)
}

