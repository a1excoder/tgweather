package main

import (
	"encoding/json"
	"os"
)

type ConfigFile struct {
	TelegramBotToken string `json:"telegram_bot_token"`
	OwmToken         string `json:"owm_token"`
}

func GetConfigFileData(fileName string) (*ConfigFile, error) {
	confFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	confBuff := make([]byte, 256)
	n, err := confFile.Read(confBuff)
	if err != nil {
		return nil, err
	}

	if err = confFile.Close(); err != nil {
		return nil, err
	}

	_data := ConfigFile{}
	if err = json.Unmarshal(confBuff[:n], &_data); err != nil {
		return nil, err
	}

	return &_data, nil
}
