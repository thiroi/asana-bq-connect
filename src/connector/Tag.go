package connector

import (
	"encoding/json"
	"golang.org/x/net/context"
)

const (
	GET_TAG_URL = "https://app.asana.com/api/1.0/tags"
)

type Tag struct {
	Id          string      `json:"gid,omitempty"`
	Name        string     `json:"name,omitempty"`
}

func loadTags(ctx context.Context)([]Tag, error){
	body, loadErr := loadAsana(ctx, GET_TAG_URL)
	if loadErr != nil {
		return nil, loadErr
	}
	tags, parseErr := parseBlobToTag(body)
	if parseErr != nil {
		return nil, parseErr
	}
	return tags, nil
}

type tagWrap struct {
	Tag []Tag `json:"data"`
}

func parseBlobToTag(blob []byte) ([]Tag, error) {
	tw := new(tagWrap)
	if err := json.Unmarshal(blob, tw); err != nil {
		return nil, err
	}

	return tw.Tag, nil
}