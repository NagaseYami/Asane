package services

import "github.com/tidwall/gjson"

type IBot interface {
	SendNasaAPOD(string, string, gjson.Result, error)
	SendEcho(string, string)
}

type Message struct {
	MessageID string
	UserID    string
	GroupID   string
	Texts     []string
	Images    []string
}

type Session struct {
	id       string
	Bot      IBot
	Messages []Message
}

func EchoMode(bot IBot, groupID string, rawMessage string) {
	if Echo(groupID, rawMessage) {
		bot.SendEcho(groupID, rawMessage)
		return
	}
}

func CommandMode(bot IBot, msg Message) {

	if len(msg.Texts) > 0 && msg.Texts[0] == "apod" {
		result, err := APOD(msg.Texts[1:])
		bot.SendNasaAPOD(msg.UserID, msg.GroupID, result, err)
		return
	}

}
