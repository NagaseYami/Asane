package qq

type LifecycleObject struct {
	MetaEventType string `json:"meta_event_type"`
	PostType      string `json:"post_type"`
	SelfID        int64  `json:"self_id"`
	SubType       string `json:"sub_type"`
	Time          int64  `json:"time"`
}

type SendPrivateMessageObject struct {
	Action string `json:"action"`
	Params struct {
		UserID     int64  `json:"user_id"`
		Message    string `json:"message"`
		AutoEscape bool   `json:"auto_escape"`
	} `json:"params"`
	Echo string `json:"echo"`
}

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

type SendGroupMessageObject struct {
	Action string `json:"action"`
	Params struct {
		GroupID    int64  `json:"group_id"`
		Message    string `json:"message"`
		AutoEscape bool   `json:"auto_escape"`
	} `json:"params"`
	Echo string `json:"echo"`
}

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
