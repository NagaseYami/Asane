package command

import (
	"Asane/internal/qq"
	"strings"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	delimiter   = " "
	baseCommand = "asane"
)

var commands = map[string]func([]string) string{
	"eroe": yandereRandomR18Illust,
	"tag":  yandereSerchTags,
}

func Excute(conn *websocket.Conn, msg qq.IReciveMessageObject) {
	slice := strings.Split(msg.GetRaw(), delimiter)
	if len(slice) < 2 || slice[0] != baseCommand {
		return
	}
	if fn, ok := commands[slice[1]]; ok {
		log.Infof("收到命令：%s", strings.Join(slice[1:], " "))

		resp := []byte{}
		if len(slice) > 2 {
			resp = msg.GetResponse(fn(slice[2:]))
		} else {
			resp = msg.GetResponse(fn([]string{}))
		}

		qq.WebSocketServer.WriteTextMessage(conn, resp)
	}
}
