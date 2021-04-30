package system

import (
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
)

const (
	BotCommand        = "asane"
	configFilePath    = "config.json"
	defaultConfigJson = `{
    "log_level": "info",
    "qq_config": {
        "enable": false,
        "address": "",
        "token": "",
		"debug": {
			"enable": false,
			"test_user_id": 0,
			"test_group_id": 0
		},
    },
    "nasa_config": {
        "enable": false,
        "api_key": ""
    },
    "echo_config": {
        "enable": false,
        "trigger_times": 2
    }
}`
)

var Config gjson.Result

func LoadConfigFile() {
	if _, err := os.Stat(configFilePath); err != nil {
		createDefaultConfigFile()
		return
	}
	b, err := os.ReadFile(configFilePath)
	HandleError(err)
	Config = gjson.ParseBytes(b)
}

func createDefaultConfigFile() {
	log.Warnln("没有检测到配置文件，生成一个空白配置文件。")
	b := []byte(defaultConfigJson)
	err := ioutil.WriteFile(configFilePath, b, 0666)
	HandleError(err)
}
