package services

import (
	"Asane/internal/qq"
	"strings"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	delimiter   = " "
	baseCommand = "Asane"
)

var commands = map[string]func([]string) string{
	"来点色图": makeMessageYandereRandomR18Illust,
}

func Excute(conn *websocket.Conn, msg qq.IReciveMessageObject) {
	slice := strings.Split(msg.GetRaw(), delimiter)
	if len(slice) < 2 || slice[0] != baseCommand {
		return
	}
	if fn, ok := commands[slice[1]]; ok {
		log.Infof("收到命令：%s", slice[1])
		resp := msg.GetResponse(fn(slice[1:]))
		qq.WebSocketServer.WriteTextMessage(conn, resp)
	}
}
