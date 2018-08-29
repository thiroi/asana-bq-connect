package main

import (
	"net/http"
	"fmt"
	"google.golang.org/appengine"
	"time"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

const MAKE_PROGRESS_QUERY = "INSERT INTO `<project>.<data_set>.task_progress` (date, userId, userName, projectId, allTask, completed, unCompleted, help, unexpected, delayed, awesome) SELECT CURRENT_TIMESTAMP() date, user.id userId, user.name userName, task.projectId projectId, COUNT(1) allTask, SUM(CASE WHEN task.completed = TRUE THEN 1 ELSE 0 END) completed, SUM(CASE WHEN task.completed = TRUE THEN 0 ELSE 1 END) uncompleted, SUM(CASE WHEN task.help = TRUE THEN 1 ELSE 0 END) help, SUM(CASE WHEN task.unexpected = TRUE THEN 1 ELSE 0 END) unexpected, SUM(CASE WHEN task.delayed = TRUE THEN 1 ELSE 0 END) delayed, SUM(CASE WHEN task.awesome = TRUE THEN 1 ELSE 0 END) awesome FROM `<project>.<data_set>.user` user INNER JOIN `<project>.<data_set>.task` task ON  user.id = task.assigneeId GROUP BY userId, userName, projectId"
const COUNT_TODAY_QUERY = "SELECT COUNT(1) AS Count FROM `<project>.<data_set>.task_progress` WHERE date > TIMESTAMP(CURRENT_DATE());"

func ProgressMaker(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Now Running!")
	basicCtx := appengine.NewContext(r)
	ctx, _ := context.WithTimeout(basicCtx, 60*time.Second)

	log.Infof(ctx,"Now initializing...")
	initConfig()
	log.Infof(ctx,"CHECKING DATA...")
	if hasDataSimple(ctx, COUNT_TODAY_QUERY) == true {
		log.Warningf(ctx,"There's already data.")
	} else {
		log.Infof(ctx,"Let's Running")
		runQuery(ctx, MAKE_PROGRESS_QUERY)
	}
	log.Infof(ctx,"ALL DONE")

}


