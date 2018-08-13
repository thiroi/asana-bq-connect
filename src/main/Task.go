package main

import (
	"golang.org/x/net/context"
	"strings"
	"strconv"
	"encoding/json"
	"time"
)

const (
	GET_TASK_WITH_SECTION_URL = "https://app.asana.com/api/1.0/sections/SECTION_ID/tasks?opt_fields=name,assignee,completed,completed_at,created_at,tags"
	SECTION_URL_KEY           = "SECTION_ID"
	DELAYED_TAG_ID            = 770623178101952
	UNEXPECTED_TAG_ID         = 770623178101955
)

type Task struct {
	ProjectId   int64
	SectionId   int64
	Id          int64
	Name        string
	AssigneeId  int64
	Delayed     bool
	Unexpected  bool
	TagIds      string
	Completed   bool
	CompletedAt time.Time
	CreatedAt   time.Time
}

type TaskJSON struct {
	Id          int64     `json:"id,omitempty"`
	Name        string    `json:"name,omitempty"`
	Completed   bool      `json:"completed,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
	Assignee    Assignee  `json:"assignee,omitempty"`
	Tags        []TaskTag `json:"tags,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
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
	jsonTagIds := taskJson.Tags
	var tagIds string
	if len(jsonTagIds) > 0 {
		tagIds = strconv.Itoa(int(jsonTagIds[0].Id))
		for i := 1; i < len(jsonTagIds); i++ {
			tagIds = tagIds + "," + strconv.Itoa(int(jsonTagIds[i].Id))
		}
	}

	return Task{
		ProjectId:   projectId,
		SectionId:   sectionId,
		Id:          taskJson.Id,
		Name:        taskJson.Name,
		AssigneeId:  taskJson.Assignee.Id,
		Delayed:     hasTagId(jsonTagIds, DELAYED_TAG_ID),
		Unexpected:  hasTagId(jsonTagIds, UNEXPECTED_TAG_ID),
		TagIds:      tagIds,
		Completed:   taskJson.Completed,
		CompletedAt: taskJson.CompletedAt,
		CreatedAt:   taskJson.CreatedAt,
	}
}

func hasTagId(tags []TaskTag, tagId int64) bool {
	for i := 0; i < len(tags); i++ {
		if tags[i].Id == tagId {
			return true
		}
	}
	return false
}
