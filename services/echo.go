package services

import (
	"github.com/NagaseYami/asane/system"
)

var groupMessageStacks = make(map[string]*messageStack)

type messageStack struct {
	messages []string
	lock     bool
}

func Echo(groupID string, rawMessage string) bool {
	// 非群消息，不做任何处理
	if groupID == "" {
		return false
	}

	// 初始化
	if _, ok := groupMessageStacks[groupID]; !ok {
		ms := &messageStack{}
		triggerTimes := int(system.Config.Get("echo_config.trigger_times").Int())
		for i := 0; i < triggerTimes; i++ {
			ms.messages = append(ms.messages, "")
		}
		groupMessageStacks[groupID] = ms
	}

	var ms = groupMessageStacks[groupID]

	// 如果已经复读过一次了，检查复读是否已经结束
	if ms.lock && rawMessage != ms.messages[0] {
		ms.lock = false
	}

	if !ms.lock {

		for i := len(ms.messages) - 1; i > 0; i-- {
			ms.messages[i] = ms.messages[i-1]
		}
		ms.messages[0] = rawMessage

		result := true
		for i := 1; i < len(ms.messages); i++ {
			if ms.messages[i] != ms.messages[0] {
				result = false
				break
			}
		}

		if result {
			ms.lock = true
		}

		return result
	}

	return false
}
