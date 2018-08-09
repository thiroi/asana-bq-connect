package main

import (
	"fmt"
	"net/http"
	"log"
	"google.golang.org/appengine"
	"context"
	"time"
)

const(
	PROJECT_PREFIX = "Sprint"
)

func Main(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Now Running!")
	log.Println("===Start===")
	start_time := time.Now()
	// initialization
	log.Println("Now initializing...")
	initConfig()
	log.Println("INITIALIZED!!!")

	// data load
	log.Println("Now Data Loading...")
	ctx := appengine.NewContext(r)
	projects, sections, tasks, tags, users, loadErr := load(ctx)
	if loadErr != nil{
		log.Print(loadErr)
	}
	log.Printf("プロジェクト数：%#v", len(projects))
	log.Printf("セクション数：%#v", len(sections))
	log.Printf("タスク数：%#v", len(tasks))
	log.Printf("タグ数：%#v", len(tags))
	log.Printf("ユーザー数：%#v", len(users))

	// data upload
	log.Println("Let's put data!")
	var bqStructs []CommonBqStruct
	bqStructs = append(bqStructs, CommonBqStruct{"project", projects})
	bqStructs = append(bqStructs, CommonBqStruct{"section", sections})
	//bqStructs = append(bqStructs, CommonBqStruct{"task", tasks})
	bqStructs = append(bqStructs, CommonBqStruct{"tag", tags})
	bqStructs = append(bqStructs, CommonBqStruct{"user", users})
	uploadErr := uploadBq(ctx, bqStructs)
	//uploadErr := putSample(ctx)
	if (uploadErr != nil){
		log.Printf("ERROR:", uploadErr)
	}
	log.Println("All done!!!")

	end_time := time.Now()
	total := end_time.Sub(start_time)
	log.Printf("TOTAL TIME:%#v", total.Seconds())
	log.Println("===End===")
}

func loadTest(ctx context.Context){
	tasksByte, loadErr := loadAsana(ctx, makeTaskUrl(770468093387339))
	if(loadErr != nil){
		log.Println(loadErr)
	}
	taskJson, parseErr := parseBlobToTaskJSON(tasksByte)
	if(parseErr != nil){
		log.Println(parseErr)
	}
	log.Println("HERE IS Tag Number")
	log.Println(len(taskJson[1].Tags))
	log.Println("HERE IS TASK JSON")
	log.Printf("%+v\n", taskJson)
}

func load(ctx context.Context)([]Project, []Section, []Task, []Tag, []User, error){
	log.Println("project loading...")
	originalProjects, projectErr := loadProjects(ctx)
	if projectErr != nil {
		return nil, nil, nil, nil, nil, projectErr
	}
	projects := filterProject(PROJECT_PREFIX, originalProjects)

	log.Println("section loading...")
	sections, sectionErr := loadSectionsWithProjects(ctx, projects)
	if sectionErr != nil {
		return nil, nil, nil, nil, nil, sectionErr
	}

	//log.Println("task loading...")
	//tasks, taskErr := loadTasksWithSections(ctx, sections)
	//if taskErr != nil {
	//	return nil, nil, nil, nil, nil, taskErr
	//}
	//
	//return projects, sections, tasks, nil, nil, nil

	log.Println("tag loading...")
	tags, tagErr := loadTags(ctx)
	if tagErr != nil {
		return nil, nil, nil, nil, nil, tagErr
	}

	log.Println("user loading...")
	users, userErr := loadUsers(ctx)
	if userErr != nil {
		return nil, nil, nil, nil, nil, userErr
	}

	return projects, sections, nil, tags, users, nil
}