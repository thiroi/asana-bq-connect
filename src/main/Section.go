package main

import (
	"time"
	"encoding/json"
	"golang.org/x/net/context"
	"strings"
	"strconv"
)

const (
	GET_SECTION_WITH_PROJECT_URL = "https://app.asana.com/api/1.0/projects/PROJECT_ID/sections?opt_fields=name,created_at"
	PROJECT_URL_KEY              = "PROJECT_ID"
)

type Section struct {
	ProjectId  int64
	Id         int64
	StoryPoint int64
	Name       string
	CreatedAt  time.Time
}

type SectionJSON struct {
	Id        int64     `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type sectionJSONWrap struct {
	SectionJSONs []SectionJSON `json:"data"`
}

func loadSectionsWithProjects(ctx context.Context, project Project) ([]Section, error) {
	var sections []Section
	sectionsByte, loadErr := loadAsana(ctx, makeSectionUrl(project.Id))
	if loadErr != nil {
		return nil, loadErr
	}
	wk, parseErr := parseBlobToSectionWithProjectId(project.Id, sectionsByte)
	if parseErr != nil {
		return nil, parseErr
	}
	sections = append(sections, wk...)
	return sections, nil
}

func makeSectionUrl(projectId int64) (string) {
	return strings.Replace(GET_SECTION_WITH_PROJECT_URL, PROJECT_URL_KEY, strconv.Itoa(int(projectId)), -1)
}

func parseBlobToSectionWithProjectId(projectId int64, blob []byte) ([]Section, error) {
	secJsons, err := parseBlobToSectionJSON(blob)
	if err != nil {
		return nil, err
	}

	var sections []Section
	for i := 0; i < len(secJsons); i++ {
		wk := convertSection(projectId, secJsons[i])
		sections = append(sections, wk)
	}
	return sections, nil
}

func parseBlobToSectionJSON(blob []byte) ([]SectionJSON, error) {
	swj := new(sectionJSONWrap)
	if err := json.Unmarshal(blob, swj); err != nil {
		return nil, err
	}
	return swj.SectionJSONs, nil
}

func convertSection(projectId int64, secJson SectionJSON) (Section) {
	name, point := splitNameAndPoint(secJson.Name)
	return Section{
		ProjectId:  projectId,
		Id:         secJson.Id,
		StoryPoint: point,
		Name:       name,
		CreatedAt:  secJson.CreatedAt,
	}
}

func splitNameAndPoint(originalName string) (string, int64) {
	var lastIndex int
	lastIndex = strings.LastIndex(originalName, " ")
	if (lastIndex < 1) {
		return originalName, 0
	}
	storyPointStr := originalName[lastIndex+1 : len(originalName)]
	storyPoint, parseErr := strconv.ParseInt(storyPointStr, 10, 32)
	if (parseErr != nil) {
		return originalName, 0
	}
	return originalName[0:lastIndex], storyPoint
}
