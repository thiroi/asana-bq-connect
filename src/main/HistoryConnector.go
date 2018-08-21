package main

import (
	"net/http"
	"os"
)

func HistoryConnector(w http.ResponseWriter, r *http.Request) {
	u, _ := r.URL.Parse(r.URL.String())
	params := u.Query()
	prefix := params.Get("sprint")
	if(prefix == ""){
		os.Exit(ERROR_INVALID_SPRINT_PARAM)
	}
	connect(w, r, prefix, true)
}

