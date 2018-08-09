package main

import (
	"log"
	"github.com/BurntSushi/toml"
)

//Config 設定ファイル
type Config struct {
	Asana AsanaConfig
	Bq BqConfig
}

type AsanaConfig struct {
	Token string
}

type BqConfig struct {
	Project string
	Dataset string
}

var config Config

func initConfig() {

	_, err := toml.DecodeFile("./config.toml", &config)
	if err != nil {
		log.Println(err)
	}
}

