package main

import (
	. "Asane/internal/qq"
	. "Asane/internal/yandere"
	"encoding/json"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func main() {
	log.Infoln("幾重にも辛酸を舐め、七難八苦を超え、艱難辛苦の果て、満願成就に至る。")

	addr := os.Getenv("WEBSOCKET_SERVER_ADDR")
	if addr == "" {
		log.Infoln("缺少环境变量WEBSOCKET_SERVER_ADDR")
		return
	}

	WebSocketServer.HandleGroupMessageFunc("来点色图", sendRandomR18IllustToGroup)
	WebSocketServer.HandlePrivateMessageFunc("来点色图", sendRandomR18IllustToPrivate)
	WebSocketServer.Run(addr, os.Getenv("WEBSOCKET_SERVER_TOKEN"))
}

func sendRandomR18IllustToGroup(result gjson.Result) {
	reciveMsg := &ReciveGroupMessageObject{}
	json.Unmarshal([]byte(result.Raw), reciveMsg)

	log.Infof("收到来自群%d的用户%d（%s）的随机色图请求", reciveMsg.GroupID, reciveMsg.Sender.UserID, reciveMsg.Sender.Nickname)

	fileURL := YandereClient.GetRandomExplicitPost().FileURL

	sendMsg := SendGroupMessageObject{}
	sendMsg.Action = "send_group_msg"
	sendMsg.Params.GroupID = reciveMsg.GroupID
	sendMsg.Params.Message = fmt.Sprintf("[CQ:image,file=%s]", fileURL)

	WebSocketServer.SendJSON(sendMsg)
}

func sendRandomR18IllustToPrivate(result gjson.Result) {
	reciveMsg := &RecivePrivateMessageObject{}
	json.Unmarshal([]byte(result.Raw), reciveMsg)

	log.Infof("收到来自用户%d的随机色图请求", reciveMsg.UserID)

	fileURL := YandereClient.GetRandomExplicitPost().FileURL

	sendMsg := SendPrivateMessageObject{}
	sendMsg.Action = "send_private_msg"
	sendMsg.Params.UserID = reciveMsg.UserID
	sendMsg.Params.Message = fmt.Sprintf("[CQ:image,file=%s]", fileURL)

	WebSocketServer.SendJSON(sendMsg)
}
