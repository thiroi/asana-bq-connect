package main

import (
	"log"
	"net/http"
	"fmt"
	"time"
	"google.golang.org/appengine"
	"os"
	"golang.org/x/net/context"
)

func connect(w http.ResponseWriter, r *http.Request, projectPrefix string, history bool) {
	fmt.Fprint(w, "Now Running!")
	log.Println("===Start===")
	start_time := time.Now()
	// initialization
	log.Println("Now initializing...")
	initConfig()
	ctx := appengine.NewContext(r)
	log.Println("INITIALIZED!!!")

	// data load
	log.Println("Now Data Loading...")
	project, sections, tasks, tags, users, loadErr := load(ctx, projectPrefix)
	if loadErr != nil {
		log.Print(loadErr)
		os.Exit(ERROR_LOADING)
	}
	log.Printf("プロジェクト数：%#v", 1)
	log.Printf("セクション数：%#v", len(sections))
	log.Printf("タスク数：%#v", len(tasks))
	log.Printf("タグ数：%#v", len(tags))
	log.Printf("ユーザー数：%#v", len(users))

	// data upload
	log.Println("Let's put data!")
	var bqStructs []CommonBqStruct
	if (history == true) {
		if hasData(ctx, "project_history", projectPrefix) == true {
			log.Println("BACKUP ERROR")
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
		log.Printf("ERROR:", uploadErr)
		os.Exit(ERROR_UPLOADING)
	}
	log.Println("All done!!!")

	end_time := time.Now()
	total := end_time.Sub(start_time)
	log.Printf("TOTAL TIME:%#v", total.Seconds())
	log.Println("===End===")
}

//func loadTest(ctx context.Context){
//	tasksByte, loadErr := loadAsana(ctx, makeTaskUrl(770468093387339))
//	if(loadErr != nil){
//		log.Println(loadErr)
//	}
//	taskJson, parseErr := parseBlobToTaskJSON(tasksByte)
//	if(parseErr != nil){
//		log.Println(parseErr)
//	}
//	log.Println("HERE IS Tag Number")
//	log.Println(len(taskJson[1].Tags))
//	log.Println("HERE IS TASK JSON")
//	log.Printf("%+v\n", taskJson)
//}

func load(ctx context.Context, projectFilter string) (Project, []Section, []Task, []Tag, []User, error) {
	log.Println("project loading...")
	originalProjects, projectErr := loadProjects(ctx)
	if projectErr != nil {
		return Project{}, nil, nil, nil, nil, projectErr
	}
	//GAEの制限上、大量のプロジェクトをターゲットにすると死ぬので制御
	project := filterProject(projectFilter, originalProjects)

	log.Println("section loading...")
	sections, sectionErr := loadSectionsWithProjects(ctx, project)
	if sectionErr != nil {
		return Project{}, nil, nil, nil, nil, sectionErr
	}

	log.Println("task loading...")
	tasks, taskErr := loadTasksWithSections(ctx, sections)
	if taskErr != nil {
		return Project{}, nil, nil, nil, nil, taskErr
	}

	log.Println("tag loading...")
	tags, tagErr := loadTags(ctx)
	if tagErr != nil {
		return Project{}, nil, nil, nil, nil, tagErr
	}

	log.Println("user loading...")
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
