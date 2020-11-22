package services

import (
	"Asane/internal/qq"
	"Asane/internal/yandere"
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

const (
	delimiter   = " "
	baseCommand = "Asane"
)

var commands = map[string]func([]string) string{
	"来点色图": makeMessageYandereRandomR18Image,
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

func makeMessageYandereRandomR18Image(params []string) string {
	post := yandere.Client.GetRandomExplicitPost()
	var imageURL = ""
	if post.FileSize < 6291456 {
		imageURL = post.FileURL
	} else if post.JpegFileSize > 0 && post.JpegFileSize < 6291456 {
		imageURL = post.JpegURL
	} else if post.SampleFileSize > 0 && post.SampleFileSize < 6291456 {
		imageURL = post.SampleURL
	} else {
		log.Debug("图太大了，换一张")
		return makeMessageYandereRandomR18Image(params)
	}
	return fmt.Sprintf("https://yande.re/post/show/%d [CQ:image,file=%s]", post.ID, imageURL)
}
