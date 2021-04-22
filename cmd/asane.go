package main

import (
	"Asane/internal/front_end/qq"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Print("幾重にも辛酸を舐め、七難八苦を超え、艱難辛苦の果て、満願成就に至る。")

	switch os.Getenv("LOG_LEVEL") {
	case "Panic":
		log.SetLevel(log.PanicLevel)
	case "Fatal":
		log.SetLevel(log.FatalLevel)
	case "Error":
		log.SetLevel(log.ErrorLevel)
	case "Warn":
		log.SetLevel(log.WarnLevel)
	case "Info":
		log.SetLevel(log.InfoLevel)
	case "Debug":
		log.SetLevel(log.DebugLevel)
	case "Trace":
		log.SetLevel(log.TraceLevel)
	}

	addr := os.Getenv("WEBSOCKET_SERVER_ADDR")
	if addr == "" {
		log.Fatal("缺少环境变量WEBSOCKET_SERVER_ADDR")
	}

	qq.WebSocketServer.HandleMessage(Excute)
	qq.WebSocketServer.Run(addr, os.Getenv("WEBSOCKET_SERVER_TOKEN"))
}
