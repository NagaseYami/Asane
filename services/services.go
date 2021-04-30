package services

import "github.com/tidwall/gjson"

type IBot interface {
	SendNasaAPOD(string, string, gjson.Result, error)
	SendEcho(string, string)
}

type Message struct {
	MessageID  string
	RawMessage string
	UserID     string
	GroupID    string
	Texts      []string
	Images     []string
}

type Session struct {
	id       string
	Bot      IBot
	Messages []Message
}

func OnReceiveMessage(bot IBot, msg Message) {

	// 复读检测
	if Echo(msg.GroupID, msg.RawMessage) {
		bot.SendEcho(msg.GroupID, msg.RawMessage)
		return
	}

	if msg.Texts[0] == "apod" {
		result, err := APOD(msg.Texts[1:])
		bot.SendNasaAPOD(msg.UserID, msg.GroupID, result, err)
		return
	}

}
