package main

import (
	"net/http"
	"fmt"
	"google.golang.org/appengine"
	"time"
	"google.golang.org/appengine/log"
	"golang.org/x/net/context"
	"os"
	"strconv"
)

func ProgressNotifier(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Now Running!")
	basicCtx := appengine.NewContext(r)
	ctx, _ := context.WithTimeout(basicCtx, 60*time.Second)

	log.Infof(ctx, "Now initializing...")
	initConfig()
	log.Infof(ctx, "CHECKING DATA...")

	data, loadErr := loadDayDiff(ctx)
	if loadErr != nil {
		log.Errorf(ctx, "Failed to load daydiff:", loadErr)
		os.Exit(ERROR_LOADING)
	}

	// convertMessage
	message := convertDataToMessage(ctx, data)

	// sendMessage
	log.Infof(ctx, message)

	sendNlope(ctx, message)
}

func convertDataToMessage(ctx context.Context, dayDiffs []DayDiff)(string){
	var message string
	for i := 0; i < len(dayDiffs); i++ {
		dayDiff := dayDiffs[i]
		log.Infof(ctx, "TEST:", dayDiff.TodayCompleted)
		log.Infof(ctx, "TEST2:", int(dayDiff.TodayCompleted))
		log.Infof(ctx, "TEST3:", strconv.Itoa(int(dayDiff.TodayCompleted)))
		if dayDiff.YesterdayCompleted == 0 & dayDiff.YesterdayUnCompleted {
			message = message + dayDiff.UserName +
				"  完了:" + strconv.Itoa(int(dayDiff.TodayCompleted)) +
				"  未完了:" + strconv.Itoa(int(dayDiff.TodayUnCompleted)) +
				"\n"
		} else {
			message = message + dayDiff.UserName +
				"  完了:" + strconv.Itoa(int(dayDiff.YesterdayCompleted)) + " -> " + strconv.Itoa(int(dayDiff.TodayCompleted)) +
				"  未完了:" + strconv.Itoa(int(dayDiff.YesterdayUnCompleted)) + " -> " + strconv.Itoa(int(dayDiff.TodayUnCompleted)) +
				"\n"
		}
	}
	return message
}
