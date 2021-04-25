package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/NagaseYami/asane/module/nasa"
	"github.com/NagaseYami/asane/module/yandere"
	"github.com/NagaseYami/asane/server/qq"
	"github.com/NagaseYami/asane/system"

	"github.com/gorilla/websocket"

	log "github.com/sirupsen/logrus"
)

const (
	baseCommand = "asane"
)

var commands = map[string]func([]string) string{
	"eroe":   yandere.YandereRandomExplicitIllust,
	"illust": yandere.YandereRandomSafeIllust,
	"tag":    yandere.YandereSerchTags,
	"apod":   nasa.NasaAPOD,
}

func main() {
	log.Print("「の人生最大の願望は、兄と兄の彼女と3Pすることだから……」")

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
	case "trace":
		log.SetLevel(log.TraceLevel)
	}

	if system.Config.QQConfig.Enable {
		qq.WebSocketServer.HandleMessage(Excute)
		qq.WebSocketServer.Run(system.Config.QQConfig.Address, system.Config.QQConfig.Token)
	}

	log.Print("「兄と添い寝、兄と添い寝……！文乃にはしない妹の特権はい勝ちぃ格付け完了……！」")
	bufio.NewScanner(os.Stdin).Scan()
}

// xcute 尝试执行QQ消息中含有的命令
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
