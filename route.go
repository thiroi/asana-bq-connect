package main

import (
	"net/http"
	"github.com/asana-connector/src/connector"
	"fmt"
	"log"
	"os"
)

func main() {
	http.HandleFunc("/init", connector.DailyInitializer)
	http.HandleFunc("/connect", connector.DailyConnector)
	http.HandleFunc("/connectHistory", connector.HistoryConnector)
	http.HandleFunc("/initHistory", connector.HistoryInitializer)
	http.HandleFunc("/makeProgress", connector.ProgressMaker)
	http.HandleFunc("/notify", connector.ProgressNotifier)
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

//func indexHandler(w http.ResponseWriter, r *http.Request) {
//	if r.URL.Path != "/" {
//		http.NotFound(w, r)
//		return
//	}
//	fmt.Fprint(w, "Hello, World!")
//}