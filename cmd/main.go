package main

import (
	"Asane/internal/qq"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Infoln("幾重にも辛酸を舐め、七難八苦を超え、艱難辛苦の果て、満願成就に至る。")

	addr := os.Getenv("WEBSOCKET_SERVER_ADDR")
	if addr == "" {
		log.Infoln("缺少环境变量WEBSOCKET_SERVER_ADDR")
		return
	}

	qq.WebSocketServer.Run(addr, os.Getenv("WEBSOCKET_SERVER_TOKEN"))
}
