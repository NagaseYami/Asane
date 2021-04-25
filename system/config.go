package system

import (
	"encoding/json"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

type ConfigFile struct {
	BaseCommands []string `json:"base_command"`
	LogLevel     string   `json:"log_level"`
	QQConfig     struct {
		Enable                    bool   `json:"enable"`
		Address                   string `json:"address"`
		Token                     string `json:"token"`
	}
	NasaConfig struct {
		Enable bool   `json:"enable"`
		APIKey string `json:"api_key"`
	}
}

const (
	configFilePath = "config.json"
)

var Config = &ConfigFile{
	BaseCommands: []string{"asane"},
	LogLevel:     "info",
}

func LoadConfigFile() {
	if _, err := os.Stat(configFilePath); err != nil {
		createDefaultConfigFile()
		return
	}
	b, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b, Config)

	if err != nil {
		log.Fatal(err)
	}
}

func createDefaultConfigFile() {
	log.Warnln("没有检测到配置文件，生成一个空白配置文件。")
	b, err := json.Marshal(Config)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(configFilePath, b, 0666)
	if err != nil {
		log.Panic(err)
	}
}
