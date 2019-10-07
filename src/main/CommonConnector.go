package main

import (
	"google.golang.org/appengine/log"
	"net/http"
	"fmt"
	"time"
	"google.golang.org/appengine"
	"os"
	"golang.org/x/net/context"
)

func connect(w http.ResponseWriter, r *http.Request, projectPrefix string, history bool) {
	fmt.Fprint(w, "Now Running!")
	basicCtx := appengine.NewContext(r)
	ctx, _ := context.WithTimeout(basicCtx, 60*time.Second)
	log.Infof(ctx, "===Start===")
	start_time := time.Now()
	// initialization
	log.Infof(ctx,"Now initializing...")
	initConfig()
	log.Infof(ctx,"INITIALIZED!!!")

	// data load
	log.Infof(ctx,"Now Data Loading...")
	project, sections, tasks, tags, users, loadErr := load(ctx, projectPrefix)
	if loadErr != nil {
		log.Errorf(ctx, loadErr.Error())
		os.Exit(ERROR_LOADING)
	}
	log.Infof(ctx,"プロジェクト数：%#v", 1)
	log.Infof(ctx,"セクション数：%#v", len(sections))
	log.Infof(ctx,"タスク数：%#v", len(tasks))
	log.Infof(ctx,"タグ数：%#v", len(tags))
	log.Infof(ctx,"ユーザー数：%#v", len(users))

	// data upload
	log.Infof(ctx,"Let's put data!")
	var bqStructs []CommonBqStruct
	if (history == true) {
		if hasData(ctx, "project_history", projectPrefix) == true {
			log.Infof(ctx,"BACKUP ERROR")
			os.Exit(ERROR_BACKUP)
		}
		bqStructs = append(bqStructs, CommonBqStruct{"project_history", project})
		bqStructs = append(bqStructs, CommonBqStruct{"section_history", sections})
		bqStructs = append(bqStructs, CommonBqStruct{"task_history", tasks})
		bqStructs = append(bqStructs, CommonBqStruct{"tag_history", tags})
		bqStructs = append(bqStructs, CommonBqStruct{"user_history", users})
	} else {
		bqStructs = append(bqStructs, CommonBqStruct{"project", project})
		bqStructs = append(bqStructs, CommonBqStruct{"section", sections})
		bqStructs = append(bqStructs, CommonBqStruct{"task", tasks})
		bqStructs = append(bqStructs, CommonBqStruct{"tag", tags})
		bqStructs = append(bqStructs, CommonBqStruct{"user", users})
	}
	uploadErr := uploadBq(ctx, bqStructs)
	//uploadErr := putSample(ctx)
	if (uploadErr != nil) {
		log.Errorf(ctx,"ERROR:", uploadErr)
		os.Exit(ERROR_UPLOADING)
	}
	log.Infof(ctx,"All done!!!")

	end_time := time.Now()
	total := end_time.Sub(start_time)
	log.Infof(ctx,"TOTAL TIME:%#v", total.Seconds())
	log.Infof(ctx,"===End===")
}

func load(ctx context.Context, projectFilter string) (Project, []Section, []Task, []Tag, []User, error) {
	log.Infof(ctx,"project loading...")
	originalProjects, projectErr := loadProjects(ctx)
	if projectErr != nil {
		return Project{}, nil, nil, nil, nil, projectErr
	}
	//GAEの制限上、大量のプロジェクトをターゲットにすると死ぬので制御
	project := filterProject(projectFilter, originalProjects)

	log.Infof(ctx,"section loading...")
	sections, sectionErr := loadSectionsWithProjects(ctx, project)
	if sectionErr != nil {
		return Project{}, nil, nil, nil, nil, sectionErr
	}

	log.Infof(ctx,"task loading...")
	tasks, taskErr := loadTasksWithSections(ctx, sections)
	if taskErr != nil {
		return Project{}, nil, nil, nil, nil, taskErr
	}

	log.Infof(ctx,"tag loading...")
	tags, tagErr := loadTags(ctx)
	if tagErr != nil {
		return Project{}, nil, nil, nil, nil, tagErr
	}

	log.Infof(ctx,"user loading...")
	users, userErr := loadUsers(ctx)
	if userErr != nil {
		return Project{}, nil, nil, nil, nil, userErr
	}

	return project, sections, tasks, tags, users, nil
}

func checkBkData(ctx context.Context, projectPrefix string) bool {
	if hasData(ctx, "project_bk", projectPrefix) == true {
		return false
	}
	return true
}
