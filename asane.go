package main

import (
	"github.com/NagaseYami/asane/bots/qq"
	"github.com/NagaseYami/asane/system"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Print("「妹の人生最大の願望は、兄と兄の彼女と3Pすることだから……」")

	system.LoadConfigFile()

	switch system.Config.Get("log_level").String() {
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "trace":
		log.SetLevel(log.TraceLevel)
	}

	if system.Config.Get("qq_config.enable").Bool() {
		qq.Server.Run()
	}

	log.Print("「兄と添い寝、兄と添い寝……！文乃にはしない妹の特権はい勝ちぃ格付け完了……！」")
}
