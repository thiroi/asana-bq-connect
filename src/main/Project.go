package main

import (
	"encoding/json"
	"golang.org/x/net/context"
	"strings"
	"time"
)

const (
	GET_PROJECT_URL = "https://app.asana.com/api/1.0/projects?opt_fields=name,archived,public,created_at"
)

type Project struct {
	Id string  `json:"gid,omitempty"`
	Name string `json:"name,omitempty"`
	Archived bool `json:"archived,omitempty"`
	Public bool `json:"public,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func filterProject(prefix string, projects []Project)(Project){
	var filteredProject = Project{}
	for i := 0; i < len(projects); i++ {
		project := projects[i]
		if strings.Contains(project.Name, prefix){
			//PREFIXを満たしており、かつ最新のものだけをとってくる
			if (filteredProject == Project{}){
				filteredProject = project
			} else if (project.CreatedAt.After(filteredProject.CreatedAt)){
				filteredProject = project
			}
		}
	}
	return filteredProject
}

func loadProjects(ctx context.Context)([]Project, error){
	body, loadErr := loadAsana(ctx, GET_PROJECT_URL)
	if loadErr != nil {
		return nil, loadErr
	}
	projects, parseErr := parseOutProjectFromData(body)
	if parseErr != nil {
		return nil, parseErr
	}
	return projects, nil
}

type projectWrap struct {
	Project []Project `json:"data"`
}

func parseOutProjectFromData(blob []byte) ([]Project, error) {
	pwj := new(projectWrap)
	if err := json.Unmarshal(blob, pwj); err != nil {
		return nil, err
	}

	return pwj.Project, nil
}
