package main

import (
	"net/http"
	"fmt"
	"log"
	"time"
	"google.golang.org/appengine"
	"os"
)

func DailyInitializer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Now Running!")
	log.Println("===Start===")
	start_time := time.Now()
	// initialization
	log.Println("Now initializing...")
	initConfig()
	ctx := appengine.NewContext(r)
	initTableErr := deleteAndCreateBq(
		ctx,
		[]CommonBqStruct{
			{"project", Project{}},
			{"section", Section{}},
			{"task", Task{}},
			{"tag", Tag{}},
			{"user", User{}},
		})
	if(initTableErr != nil){
		log.Printf("ERROR:", initTableErr)
		os.Exit(ERROR_DELETING)
	}
	log.Println("INITIALIZED!!!")
	end_time := time.Now()
	total := end_time.Sub(start_time)
	log.Printf("TOTAL TIME:%#v", total.Seconds())
	log.Println("===End===")
}
