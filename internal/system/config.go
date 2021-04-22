package system

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

type ConfigFile struct {
	LogLevel           string `json:"log_level"`
	MediaDirectoryPath string `json:"media_directory_path"`
	QQConfig           struct {
		Enable            bool   `json:"enable"`
		Address           string `json:"address"`
		Token             string `json:"token"`
		DataDirectoryPath string `json:"data_directory_path"`
	}
	YandereConfig struct {
		Enable bool `json:"enable"`
	}
	NasaConfig struct {
		Enable bool   `json:"enable"`
		APIKey string `json:"api_key"`
	}
}

const (
	configFilePath = "config.json"
	imageDirectory = "Images"
)

var Config = &ConfigFile{
	LogLevel:           "info",
	MediaDirectoryPath: "Media",
}

var ImageDirectoryPath string

func LoadConfigFile() {
	if _, err := os.Stat(configFilePath); err != nil {
		createDefaultConfigFile()
	}
	b, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(b, Config)

	if err != nil {
		log.Fatalln(err)
	}

	ImageDirectoryPath = path.Join(Config.MediaDirectoryPath, imageDirectory)
}

func createDefaultConfigFile() {
	log.Println("没有检测到配置文件，生成一个空白配置文件")
	b, err := json.Marshal(Config)
	if err != nil {
		log.Panicln(err)
	}
	ioutil.WriteFile(configFilePath, b, 0666)
}

func CreateDirectory() {
	if Config.MediaDirectoryPath != "" {
		os.MkdirAll(ImageDirectoryPath, 0666)
	}
}
