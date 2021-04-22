package main

import (
	"Asane/internal/server/qq"
	"Asane/internal/system"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Print("「妹の人生最大の願望は、兄と兄の彼女と3Pすることだから……」")

	system.LoadConfigFile()
	system.CreateDirectory()

	switch system.Config.LogLevel {
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

	if system.Config.QQConfig.Enable {
		qq.WebSocketServer.HandleMessage(Excute)
		qq.WebSocketServer.Run(system.Config.QQConfig.Address, system.Config.QQConfig.Token)
	}

	log.Print("「兄と添い寝、兄と添い寝……！文乃にはしない妹の特権はい勝ちぃ格付け完了……！」")
}
