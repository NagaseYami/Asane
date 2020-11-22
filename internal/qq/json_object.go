package qq

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type IReciveMessageObject interface {
	GetSelfID() int64
	GetRaw() string
	GetResponse(string) []byte
}

// LifecycleObject 元事件：生命周期
type LifecycleObject struct {
	MetaEventType string `json:"meta_event_type"`
	PostType      string `json:"post_type"`
	SelfID        int64  `json:"self_id"`
	SubType       string `json:"sub_type"`
	Time          int64  `json:"time"`
}

// SendPrivateMessageObject 发送私聊消息
type SendPrivateMessageObject struct {
	Action string `json:"action"`
	Params struct {
		UserID     int64  `json:"user_id"`
		Message    string `json:"message"`
		AutoEscape bool   `json:"auto_escape"`
	} `json:"params"`
	Echo string `json:"echo"`
}

// RecivePrivateMessageObject 接收私聊消息
type RecivePrivateMessageObject struct {
	Font        int32  `json:"font"`
	Message     string `json:"message"`
	MessageID   int32  `json:"message_id"`
	MessageType string `json:"message_type"`
	PostType    string `json:"post_type"`
	RawMessage  string `json:"raw_message"`
	SelfID      int64  `json:"self_id"`
	Sender      struct {
		Age      int32  `json:"age"`
		Nickname string `json:"nickname"`
		Sex      string `json:"sex"`
		UserID   int64  `json:"user_id"`
	} `json:"sender"`
	SubType string `json:"sub_type"`
	Time    int64  `json:"time"`
	UserID  int64  `json:"user_id"`
}

func (obj *RecivePrivateMessageObject) GetSelfID() int64 {
	return obj.SelfID
}

func (obj *RecivePrivateMessageObject) GetRaw() string {
	return obj.RawMessage
}

func (obj *RecivePrivateMessageObject) GetResponse(msg string) []byte {
	sendMsg := SendPrivateMessageObject{}
	sendMsg.Action = "send_private_msg"
	sendMsg.Params.UserID = obj.UserID
	sendMsg.Params.Message = msg
	sendMsg.Params.AutoEscape = false
	result, err := json.Marshal(sendMsg)
	if err != nil {
		log.Error(err)
	}
	return result
}

// SendGroupMessageObject 发送群消息
type SendGroupMessageObject struct {
	Action string `json:"action"`
	Params struct {
		GroupID    int64  `json:"group_id"`
		Message    string `json:"message"`
		AutoEscape bool   `json:"auto_escape"`
	} `json:"params"`
	Echo string `json:"echo"`
}

// ReciveGroupMessageObject 接收群消息
type ReciveGroupMessageObject struct {
	Anonymous struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
		Flag string `json:"flag"`
	} `json:"anonymous"`
	Font        int32  `json:"font"`
	GroupID     int64  `json:"group_id"`
	Message     string `json:"message"`
	MessageID   int32  `json:"message_id"`
	MessageType string `json:"message_type"`
	PostType    string `json:"post_type"`
	RawMessage  string `json:"raw_message"`
	SelfID      int64  `json:"self_id"`
	Sender      struct {
		Age      int32  `json:"age"`
		Area     string `json:"area"`
		Card     string `json:"card"`
		Level    string `json:"level"`
		Nickname string `json:"nickname"`
		Role     string `json:"role"`
		Sex      string `json:"sex"`
		Title    string `json:"title"`
		UserID   int64  `json:"user_id"`
	} `json:"sender"`
	SubType string `json:"sub_type"`
	Time    int64  `json:"time"`
	UserID  int64  `json:"user_id"`
}

func (obj *ReciveGroupMessageObject) GetSelfID() int64 {
	return obj.SelfID
}

func (obj *ReciveGroupMessageObject) GetRaw() string {
	return obj.RawMessage
}

func (obj *ReciveGroupMessageObject) GetResponse(msg string) []byte {
	sendMsg := SendGroupMessageObject{}
	sendMsg.Action = "send_group_msg"
	sendMsg.Params.GroupID = obj.GroupID
	sendMsg.Params.Message = msg
	sendMsg.Params.AutoEscape = false
	result, err := json.Marshal(sendMsg)
	if err != nil {
		log.Error(err)
	}
	return result
}
