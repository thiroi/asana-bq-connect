package main

import (
	"encoding/json"
	"golang.org/x/net/context"
)

const (
	GET_USER_URL = "https://app.asana.com/api/1.0/users"
)

type User struct {
	Id          string      `json:"gid,omitempty"`
	Name        string     `json:"name,omitempty"`
}

func loadUsers(ctx context.Context)([]User, error){
	body, loadErr := loadAsana(ctx, GET_USER_URL)
	if loadErr != nil {
		return nil, loadErr
	}
	users, parseErr := parseBlobToUser(body)
	if parseErr != nil {
		return nil, parseErr
	}
	return users, nil
}

type userWrap struct {
	User []User `json:"data"`
}

func parseBlobToUser(blob []byte) ([]User, error) {
	tw := new(userWrap)
	if err := json.Unmarshal(blob, tw); err != nil {
		return nil, err
	}

	return tw.User, nil
}