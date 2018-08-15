package src

import (
"net/http"
"src/main"
)

func init() {
	// net/http
	http.HandleFunc("/init", main.DailyInitializer)
	http.HandleFunc("/connect", main.DailyConnector)
	http.HandleFunc("/connectHistory", main.HistoryConnector)
	http.HandleFunc("/initHistory", main.HistoryInitializer)
}
