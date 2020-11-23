package services

import (
	"Asane/internal/qq"
	"strings"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	baseCommand = "asane"
)

var commands = map[string]func([]string) string{
	"eroe": yandereRandomR18Illust,
	"tag":  yandereSerchTags,
	"apod": nasaAPOD,
}

// Excute 尝试执行QQ消息中含有的命令
func Excute(conn *websocket.Conn, msg qq.IReciveMessageObject) {
	slice := strings.Split(msg.GetRawMessage(), " ")
	if len(slice) < 2 || slice[0] != baseCommand {
		return
	}
	if fn, ok := commands[slice[1]]; ok {
		log.Infof("收到命令：%s", strings.Join(slice[1:], " "))

		resp := []byte{}
		if len(slice) > 2 {
			resp = msg.Bytes(fn(slice[2:]))
		} else {
			resp = msg.Bytes(fn([]string{}))
		}

		qq.WebSocketServer.WriteTextMessage(conn, resp)
	}
}
