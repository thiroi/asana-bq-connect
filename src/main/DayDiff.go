package main

import (
	"log"
	"google.golang.org/api/iterator"
	"golang.org/x/net/context"
)

const DAY_DIFF_SQL = "SELECT today.userId, today.userName, today.completed AS todayCompleted, today.uncompleted AS todayUnCompleted, yesterday.completed AS yesterdayCompleted, yesterday.uncompleted AS yesterdayUnCompleted FROM ( SELECT   projectId,   userId,   userName,   completed,   uncompleted FROM   `<project>.<data_set>.task_progress` WHERE   date > TIMESTAMP(CURRENT_DATE())) today LEFT JOIN ( SELECT   projectId,   userId,   userName,   completed,   uncompleted  FROM   `<project>.<data_set>.task_progress` WHERE   date BETWEEN TIMESTAMP(DATE_SUB(CURRENT_DATE(), INTERVAL 1 DAY))   AND TIMESTAMP(CURRENT_DATE())) yesterday ON today.userId = yesterday.userId AND today.projectId = yesterday.projectId;"

type DayDiff struct {
	UserId int64  `json:"id,omitempty"`
	UserName string `json:"name,omitempty"`
	TodayCompleted int64 `json:"archived,omitempty"`
	TodayUnCompleted int64 `json:"archived,omitempty"`
	YesterdayCompleted int64 `json:"archived,omitempty"`
	YesterdayUnCompleted int64 `json:"archived,omitempty"`
}


func loadDayDiff(ctx context.Context)([]DayDiff, error){
	it, err := readQuery(ctx, DAY_DIFF_SQL)
	if err != nil {
		return nil, err
	}
	var result []DayDiff
	for {
		var line DayDiff
		err := it.Next(&line)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println("Failed to Iterate Query:%v", err)
		}
		result = append(result, line)
	}
	return result, nil
}