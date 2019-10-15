package connector

import (
	"net/http"
	"fmt"
	"log"
	"time"
	"os"
)

func HistoryInitializer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Now Running!")
	ctx := r.Context()
	log.Println( "===Start===")
	start_time := time.Now()
	// initialization
	log.Println("Now initializing...")
	initConfig()
	initTableErr := deleteAndCreateBq(
		ctx,
		[]CommonBqStruct{
			{"project_history", Project{}},
			{"section_history", Section{}},
			{"task_history", Task{}},
			{"tag_history", Tag{}},
			{"user_history", User{}},
		})
	if(initTableErr != nil){
		log.Println("ERROR: ", initTableErr)
		os.Exit(ERROR_DELETING)
	}
	log.Println("INITIALIZED!!!")
	end_time := time.Now()
	total := end_time.Sub(start_time)
	log.Println("TOTAL TIME:", total.Seconds())
	log.Println("===End===")
}
