package connector

import (
	"time"
	"github.com/kokardy/jpholiday"
	"log"
)


func truncateTimeToDate(from time.Time)(time.Time){
	to := from.Truncate( time.Hour ).Add( - time.Duration(from.Hour()) * time.Hour )
	return to
}

func getSprintEndDate(projectName string)(time.Time){
	ymdStr := "20" + projectName[len(projectName)-6 : len(projectName)]
	sprintEndDay, err := time.Parse("20060102", ymdStr)
	if err != nil {
		log.Fatal("日付の取得ができませんでした。プロジェクト名を確認してください。プロジェクト名はSprint-YYmmDDである必要があります。 入力値:", ymdStr)
	}
	return sprintEndDay
}

func countBusinessDays(startDate time.Time, endDate time.Time)(int){
	targetDate := startDate
	counter := 0
	for targetDate.Before(endDate) || targetDate.Equal(endDate) {
		if(isBusinessDay(targetDate)){
			log.Println(targetDate)
			counter++
		}
		targetDate = targetDate.AddDate(0 ,0, 1)
	}
	return counter
}

func isBusinessDay(date time.Time)(bool){
	if(date.Weekday() == time.Saturday || date.Weekday() == time.Sunday){
		return false;
	}
	isHoliday, _ := jpholiday.TimeToDate(date).Holiday()
	if(isHoliday){
		return false;
	}
	return true;
}