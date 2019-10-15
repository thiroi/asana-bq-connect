package connector

import (
	"net/http"
	"fmt"
	"time"
	"log"
	"golang.org/x/net/context"
	"os"
	"strconv"
)

func ProgressNotifier(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Now Running!")
	basicCtx := r.Context()
	ctx, _ := context.WithTimeout(basicCtx, 60*time.Second)

	log.Println( "Now initializing...")
	initConfig()
	log.Println( "CHECKING DATA...")

	data, loadErr := loadDayDiff(ctx)
	if loadErr != nil {
		log.Println( "Failed to load daydiff:", loadErr)
		os.Exit(ERROR_LOADING)
	}

	reporter, reportErr := getTodayReporter(ctx)
	if reportErr != nil {
		log.Println( "Failed to load reporter:", loadErr)
		os.Exit(ERROR_LOADING)
	}

	// convertMessage
	message := convertDataToMessage(ctx, data)

	sendNlope(ctx, message, reporter)
}

func convertDataToMessage(ctx context.Context, dayDiffs []DayDiff)(string){
	var message string
	var awesomeMan string
	var bestScore int64 = 0

	for i := 0; i < len(dayDiffs); i++ {
		dayDiff := dayDiffs[i]
		if bestScore < dayDiff.TodayAwesome {
			bestScore = dayDiff.TodayAwesome
			awesomeMan = dayDiff.UserName
		}
		if dayDiff.YesterdayCompleted == 0 && dayDiff.YesterdayUnCompleted == 0 {
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
	if bestScore > 0 {
		message = message + "\n" + "\n" + "【今最もAWESOMEな人】　" + ":parrot:" + awesomeMan + " "  + ":parrot:" +
			"\n" + "現在 " + strconv.Itoa(int(bestScore))  + " AWESOMEです" + "\n"
	}
	return message
}

func getBusinessDate(projectName string)(int, int){
	endDate := getSprintEndDate(projectName)
	firstDate := endDate.AddDate(0, 0, -6)
	today := truncateTimeToDate(time.Now())

	wholeDays := countBusinessDays(firstDate, endDate)
	remainingDays := countBusinessDays(today, endDate)

	return wholeDays, remainingDays
}
