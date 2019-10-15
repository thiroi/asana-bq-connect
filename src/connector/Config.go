package connector

import (
	"log"
	"github.com/BurntSushi/toml"
)

//Config 設定ファイル
type Config struct {
	Asana AsanaConfig
	Bq BqConfig
	Slack SlackConfig
}

type AsanaConfig struct {
	Token string
}

type BqConfig struct {
	Project string
	Dataset string
}

type SlackConfig struct {
	Token string
}

var config Config

func initConfig() {

	_, err := toml.DecodeFile("./config.toml", &config)
	if err != nil {
		log.Println(err)
	}
}

