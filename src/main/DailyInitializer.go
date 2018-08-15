package main

import (
	"net/http"
	"fmt"
	"google.golang.org/appengine/log"
	"time"
	"google.golang.org/appengine"
	"os"
)

func DailyInitializer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Now Running!")
	ctx := appengine.NewContext(r)
	log.Infof(ctx, "===Start===")
	start_time := time.Now()
	// initialization
	log.Infof(ctx, "Now initializing...")
	initConfig()
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
		log.Errorf(ctx, "ERROR: %v", initTableErr)
		os.Exit(ERROR_DELETING)
	}
	log.Infof(ctx, "INITIALIZED!!!")
	end_time := time.Now()
	total := end_time.Sub(start_time)
	log.Infof(ctx, "TOTAL TIME:%#v", total.Seconds())
	log.Infof(ctx, "===End===")
}
