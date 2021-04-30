package qq

import (
	"github.com/NagaseYami/asane/services"
	"github.com/NagaseYami/asane/system"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type bot struct {
	id   string
	conn *websocket.Conn
}

func (b *bot) SendNasaAPOD(userID string, groupID string, apod gjson.Result, err error) {
	isPrivate := groupID == ""

	var msg string

	if err != nil {
		msg, _ = sjson.Set(``, "params.message.0.type", "text")
		msg, _ = sjson.Set(msg, "params.message.0.data.text", err.Error())
	} else if isPrivate {
		msg, _ = sjson.Set(``, "params.message.0.type", "text")
		msg, _ = sjson.Set(msg, "params.message.0.data.text", apod.Get("title").String()+"\n")
		if apod.Get("media_type").String() == "image" {
			msg, _ = sjson.Set(msg, "params.message.1.type", "image")
			msg, _ = sjson.Set(msg, "params.message.1.data.file", apod.Get("url").String())
		} else if apod.Get("media_type").String() == "video" {
			msg, _ = sjson.Set(msg, "params.message.1.type", "video")
			msg, _ = sjson.Set(msg, "params.message.1.data.file", apod.Get("url").String())
		}
		msg, _ = sjson.Set(msg, "params.message.2.type", "text")
		msg, _ = sjson.Set(msg, "params.message.2.data.text", "\n"+apod.Get("explanation").String())
	} else if !isPrivate {
		msg, _ = sjson.Set(``, "params.message.0.type", "at")
		msg, _ = sjson.Set(msg, "params.message.0.data.qq", userID)
		msg, _ = sjson.Set(msg, "params.message.1.type", "text")
		msg, _ = sjson.Set(msg, "params.message.1.data.text", "\n"+apod.Get("title").String()+"\n")
		if apod.Get("media_type").String() == "image" {
			msg, _ = sjson.Set(msg, "params.message.2.type", "image")
			msg, _ = sjson.Set(msg, "params.message.2.data.file", apod.Get("url").String())
		} else if apod.Get("media_type").String() == "video" {
			msg, _ = sjson.Set(msg, "params.message.2.type", "video")
			msg, _ = sjson.Set(msg, "params.message.2.data.file", apod.Get("url").String())
		}
		msg, _ = sjson.Set(msg, "params.message.3.type", "text")
		msg, _ = sjson.Set(msg, "params.message.3.data.text", "\n"+apod.Get("explanation").String())
	}

	b.sendMessage(userID, groupID, msg)

}

func (b *bot) SendEcho(groupID string, rawMessage string) {
	b.sendMessage("", groupID, rawMessage)
}

func (b *bot) universalRouter(result gjson.Result) {
	echo := result.Get("echo").String()
	if echo == "" {
		switch result.Get("post_type").String() {
		case "message":
			log.Trace(result.String())
			b.messageRouter(result)

		case "notice":

		case "request":

		case "meta_event":
			b.metaEventRouter(result)

		}
	} else {

	}
}

func (b *bot) metaEventRouter(result gjson.Result) {
	switch result.Get("meta_event_type").String() {
	case "lifecycle":

	case "heartbeat":

	}
}

func (b *bot) messageRouter(result gjson.Result) {

	var message services.Message
	var anyCall bool

	anyCall, message.Texts, message.Images = b.messageSlicer(result.Get("message"))

	switch result.Get("message_type").String() {
	case "private":
		// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#%E7%A7%81%E8%81%8A%E6%B6%88%E6%81%AF
	case "group":
		// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#%E7%BE%A4%E6%B6%88%E6%81%AF
		if anyCall {
			message.GroupID = result.Get("group_id").String()
		} else {
			return
		}
	}

	message.MessageID = result.Get("message_id").String()
	message.RawMessage = result.Get("raw_message").String()
	message.UserID = result.Get("user_id").String()

	var iBot services.IBot = b
	services.OnReceiveMessage(iBot, message)
}

func (b *bot) messageSlicer(msg gjson.Result) (bool, []string, []string) {
	// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/array.md
	at := msg.Get(`#(type=="at")#`).Array()
	anyCall := false
	for _, v := range at {
		if v.Get("data.qq").String() == b.id {
			anyCall = true
			break
		}
	}

	var texts []string
	for _, v := range msg.Get(`#(type=="text")#`).Array() {
		texts = append(texts, v.Get(`data.text`).String())
	}

	var images []string
	for _, v := range msg.Get(`#(type=="image")#`).Array() {
		images = append(images, v.Get(`data.url`).String())
	}

	return anyCall, texts, images

}

func (b *bot) sendMessage(userID string, groupID string, msg string) {
	if groupID == "" {
		msg, _ = sjson.Set(msg, "action", "send_private_msg")
		msg, _ = sjson.Set(msg, "params.user_id", userID)
	} else {
		msg, _ = sjson.Set(msg, "action", "send_group_msg")
		msg, _ = sjson.Set(msg, "params.group_id", groupID)
	}
	msg, _ = sjson.Set(msg, "echo", "")

	log.Tracef("BOT %v 发送了TextMessage至go-cqhttp", b.id)
	err := b.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	system.HandleError(err)
}
