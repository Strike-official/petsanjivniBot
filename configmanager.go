package main

import (
	"encoding/json"
	"io/ioutil"
)

type AppConfig struct {
	BaseURL          string `json:"base_url"`
	TelegramAPIToken string `json:"telegramAPIToken"`
	RDSCredentials   string `json:"rdsCredentials"`
	Port             string `json:"port"`
	MongoSourceUrl   string `json:"mongo_source_url"`
	MongoDBName      string `json:"mongo_db_name"`
	GCApi            string `json:"gc_api"`
}

var (
	AppConf *AppConfig
)

func InitAppConfig(file string) error {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(raw, &AppConf); err != nil {
		return err
	}
	return nil
}

func GetAppConfig() *AppConfig {
	return AppConf
}
