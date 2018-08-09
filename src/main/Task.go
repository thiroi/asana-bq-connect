package main

import (
	"golang.org/x/net/context"
	"strings"
	"strconv"
	"encoding/json"
	"time"
)

const (
	GET_TASK_WITH_SECTION_URL = "https://app.asana.com/api/1.0/sections/SECTION_ID/tasks?opt_fields=name,assignee,created_at,tags"
	SECTION_URL_KEY           = "SECTION_ID"
)

type Task struct {
	Id          int64
	Projectid   int64
	Sectionid   int64
	Name        string
	Completed   bool
	Completed_at time.Time
	Assigneeid  int64
	Created_at   time.Time
	Tagids      []int64
}

type TaskJSON struct {
	Id          int64      `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Completed   bool       `json:"completed,omitempty"`
	Completed_at time.Time `json:"completedAt,omitempty"`
	Assignee    Assignee   `json:"assignee,omitempty"`
	Tags        []TaskTag  `json:"tags,omitempty"`
	Created_at   time.Time `json:"created_at,omitempty"`
}

type TaskTag struct {
	Id int64 `json:"id,omitempty"`
}

type Assignee struct {
	Id int64 `json:"id,omitempty"`
}

type taskJSONWrap struct {
	TaskJSONs []TaskJSON `json:"data"`
}

func loadTasksWithSections(ctx context.Context, sections []Section) ([]Task, error) {
	var tasks []Task
	for i := 0; i < len(sections); i++ {
		section := sections[i]
		tasksByte, loadErr := loadAsana(ctx, makeTaskUrl(section.Id))
		if loadErr != nil {
			return nil, loadErr
		}
		wk, parseErr := parseBlobToTaskWithSection(section, tasksByte)
		if parseErr != nil {
			return nil, parseErr
		}
		tasks = append(tasks, wk...)
	}
	return tasks, nil
}

func makeTaskUrl(sectionId int64) (string) {
	return strings.Replace(GET_TASK_WITH_SECTION_URL, SECTION_URL_KEY, strconv.Itoa(int(sectionId)), -1)
}

func parseBlobToTaskWithSection(section Section, blob []byte) ([]Task, error) {
	taskJsons, err := parseBlobToTaskJSON(blob)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	for i := 0; i < len(taskJsons); i++ {
		wk := convertTask(section.ProjectId, section.Id, taskJsons[i])
		tasks = append(tasks, wk)
	}
	return tasks, nil
}

func parseBlobToTaskJSON(blob []byte) ([]TaskJSON, error) {
	tjw := new(taskJSONWrap)
	if err := json.Unmarshal(blob, tjw); err != nil {
		return nil, err
	}
	return tjw.TaskJSONs, nil
}

func convertTask(projectId, sectionId int64, taskJson TaskJSON) (Task) {
	var tagIds []int64
	jsonTagIds := taskJson.Tags
	for i:=0; i< len(jsonTagIds); i++ {
		tagIds = append(tagIds, jsonTagIds[i].Id)
	}

	return Task{
		Id:          taskJson.Id,
		Projectid:   projectId,
		Sectionid:   sectionId,
		Name:        taskJson.Name,
		Completed:   taskJson.Completed,
		Completed_at: taskJson.Completed_at,
		Assigneeid:  taskJson.Assignee.Id,
		Created_at:   taskJson.Created_at,
		Tagids:      tagIds,
	}
}
