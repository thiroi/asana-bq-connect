package connector

import (
	"time"
	"encoding/json"
	"golang.org/x/net/context"
	"strings"
	"strconv"
	"log"
)

const (
	GET_SECTION_WITH_PROJECT_URL = "https://app.asana.com/api/1.0/projects/PROJECT_ID/sections?opt_fields=name,created_at"
	PROJECT_URL_KEY              = "PROJECT_ID"
)

type Section struct {
	ProjectId  string
	Id         string
	StoryPoint int64
	Name       string
	CreatedAt  time.Time
}

type SectionJSON struct {
	Id        string     `json:"gid,omitempty"`
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
	wk, parseErr := parseBlobToSectionWithProjectId(ctx, project.Id, sectionsByte)
	if parseErr != nil {
		return nil, parseErr
	}
	sections = append(sections, wk...)
	return sections, nil
}

func makeSectionUrl(projectId string) (string) {
	return strings.Replace(GET_SECTION_WITH_PROJECT_URL, PROJECT_URL_KEY, projectId, -1)
}

func parseBlobToSectionWithProjectId(ctx context.Context, projectId string, blob []byte) ([]Section, error) {
	secJsons, err := parseBlobToSectionJSON(blob)
	if err != nil {
		return nil, err
	}

	var sections []Section
	for i := 0; i < len(secJsons); i++ {
		wk := convertSection(ctx, projectId, secJsons[i])
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

func convertSection(ctx context.Context, projectId string, secJson SectionJSON) (Section) {
	name, point := splitNameAndPoint(ctx, secJson.Name)
	return Section{
		ProjectId:  projectId,
		Id:         secJson.Id,
		StoryPoint: point,
		Name:       name,
		CreatedAt:  secJson.CreatedAt,
	}
}

func splitNameAndPoint(ctx context.Context, originalName string) (string, int64) {
	// Trimする
	// 一番最後の文字と、それ以外でわける
	// 数字かどうかを確認する
	// 数字であればストーリーポイント扱い、そうでなければエラーメッセージを表示

	trimedName := strings.TrimSpace(originalName)
	storyPointStr := trimedName[len(trimedName)-1:len(trimedName)]
	nameWithoutPoint := strings.TrimSpace(trimedName[0:len(trimedName)-1])
	storyPoint, parseErr := strconv.ParseInt(storyPointStr, 10, 32)
	if (parseErr != nil) {
		log.Println("ストーリーポイントが読み取れないセクションです：" + originalName)
		return nameWithoutPoint, 0
	}
	return nameWithoutPoint, storyPoint
}
