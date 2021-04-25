package services

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/NagaseYami/asane/services/nasa"
	"github.com/NagaseYami/asane/system"
)

type ISession interface {
	GetBotID() string
	GetRawMessage() string
	SendPictureSearchResult()
	SendNasaAPOD(nasa.APODResponseObject, error)
	SendHelp(string)
}

func OnRecivedMessage(session ISession) {

	valid, slice := Pretreatment(session.GetRawMessage(), session.GetBotID())

	if !valid {
		return
	} else if slice == nil {
		// FIXME nagase : 返回帮助信息
		return
	}

	command := slice[0]
	params := slice[1:]

	switch command {
	case "apod":
		session.SendNasaAPOD(nasa.APOD(params))
	}
}

func Pretreatment(raw string, botID string) (bool, []string) {
	// QQ中有人@BOT的情况
	str := fmt.Sprintf(`\[CQ:at,qq=%v\]`, botID)

	// FIXME nagase : Discord中有人@BOT的情况

	for _, value := range system.Config.BaseCommands {
		str += fmt.Sprintf("|%v", value)
	}

	re := regexp.MustCompile(`.*?` + `(` + str + `) *`)
	if re.MatchString(raw) {
		return true, SliceMessage(re.ReplaceAllString(raw, ""))
	}

	return false, nil
}

func SliceMessage(msg string) []string {
	result := strings.Split(msg, " ")
	if len(result) == 0 {
		return nil
	}
	return result
}
