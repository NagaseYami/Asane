package qq

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/NagaseYami/asane/services"
	"github.com/NagaseYami/asane/services/nasa"
	"github.com/NagaseYami/asane/system"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

var Server = &qqServer{
	token: system.Config.QQConfig.Token,
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type qqBot struct {
	ID   string
	conn *websocket.Conn
}

func (bot *qqBot) SendTextMessage(b []byte) {
	log.Traceln("BOT %v 发送了TextMessage至cqhttp", bot.ID)
	bot.conn.WriteMessage(websocket.TextMessage, b)
}

type qqSession struct {
	Bot           *qqBot
	TargetUserID  int64
	TargetGroupID int64
	RawMessage    string
}

func (s *qqSession) GetBotID() string {
	return s.Bot.ID
}

func (s *qqSession) GetRawMessage() string {
	return s.RawMessage
}

func (s *qqSession) SendPictureSearchResult() {

}

func (s *qqSession) SendNasaAPOD(obj nasa.APODResponseObject, err error) {
	var msg string

	if err != nil {
		msg = err.Error()
	} else {
		msg = fmt.Sprintf("[CQ:%v,file=%v]\n%v\n%v", obj.MediaType, obj.URL, obj.Title, obj.Explanation)
	}

	if s.TargetUserID != 0 {
		s.sendPrivateMessage(msg)
	} else if s.TargetGroupID != 0 {
		s.sendGroupMessage(msg)
	}
}

func (s *qqSession) SendHelp(str string) {

}

func (s *qqSession) sendPrivateMessage(msg string) {
	sendMsg := &SendPrivateMessageObject{}
	sendMsg.Action = "send_private_msg"
	sendMsg.Params.UserID = s.TargetUserID
	sendMsg.Params.Message = msg
	sendMsg.Params.AutoEscape = false
	bytes, err := json.Marshal(sendMsg)
	if err != nil {
		log.Panic(err)
	}
	s.Bot.conn.WriteMessage(websocket.TextMessage, bytes)
}

func (s *qqSession) sendGroupMessage(msg string) {
	sendMsg := &SendGroupMessageObject{}
	sendMsg.Action = "send_group_msg"
	sendMsg.Params.GroupID = s.TargetGroupID
	sendMsg.Params.Message = msg
	sendMsg.Params.AutoEscape = false
	bytes, err := json.Marshal(sendMsg)
	if err != nil {
		log.Panic(err)
	}
	s.Bot.conn.WriteMessage(websocket.TextMessage, bytes)
}

type qqServer struct {
	token string
	bots  map[string]*qqBot
}

func (s *qqServer) Init() {
	s.bots = make(map[string]*qqBot)
}

func (s *qqServer) Run(addr, authToken string) {
	s.token = authToken
	http.HandleFunc("/", s.any)
	log.Infof("go-cqhttp专用反向Websocket服务器已启动: %v", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (s *qqServer) any(w http.ResponseWriter, r *http.Request) {
	if !s.authorization(r) {
		log.Warnf("已拒绝 %v 的 Websocket 请求: Token鉴权失败", r.RemoteAddr)
		w.WriteHeader(401)
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warnf("处理 Websocket 请求时出现错误: %v", err)
		return
	}
	botID := r.Header.Get("X-Self-ID")
	log.Infof("接受来自QQ号：%v 的 Websocket 连接: %v (/)", botID, r.RemoteAddr)
	bot := &qqBot{
		ID:   botID,
		conn: c,
	}
	s.bots[botID] = bot
	go s.listenEvent(bot)
}

func (s *qqServer) authorization(r *http.Request) bool {
	if s.token != "" {
		if auth := r.URL.Query().Get("access_token"); auth != s.token {
			if auth := strings.SplitN(r.Header.Get("Authorization"), " ", 2); len(auth) != 2 || auth[1] != s.token {
				log.Warnf("已拒绝 %v 的 Websocket 请求: Token鉴权失败", r.RemoteAddr)
				return false
			}
		}
	}
	return true
}

func (s *qqServer) listenEvent(bot *qqBot) {
	defer bot.conn.Close()
	defer delete(s.bots, bot.ID)
	for {
		t, payload, err := bot.conn.ReadMessage()
		if err != nil {
			log.Warnf("监听反向WS API时出现错误: %v", err)
			break
		}
		if t == websocket.TextMessage {
			s.universalRouter(bot, gjson.ParseBytes(payload))
		}
	}
}

func (s *qqServer) universalRouter(bot *qqBot, result gjson.Result) {
	switch result.Get("post_type").Str {
	case "message":
		log.Trace(result.String())
		s.messageRouter(bot, result)

	case "notice":

	case "request":

	case "meta_event":
		s.metaEventRouter(bot, result)

	}
}

func (s *qqServer) metaEventRouter(bot *qqBot, result gjson.Result) {
	switch result.Get("meta_event_type").Str {
	case "lifecycle":

	case "heartbeat":

	}
}

func (s *qqServer) messageRouter(bot *qqBot, result gjson.Result) {
	switch result.Get("message_type").Str {
	case "private":
		reciveMsg := &RecivePrivateMessageObject{}
		json.Unmarshal([]byte(result.Raw), reciveMsg)

		var session qqSession
		session.Bot = bot
		session.RawMessage = reciveMsg.RawMessage
		session.TargetUserID = reciveMsg.UserID
		session.TargetGroupID = 0

		var iSession services.ISession
		iSession = &session
		services.OnRecivedMessage(iSession)

	case "group":
		reciveMsg := &ReciveGroupMessageObject{}
		json.Unmarshal([]byte(result.Raw), reciveMsg)

		var session qqSession
		session.Bot = bot
		session.RawMessage = reciveMsg.RawMessage
		session.TargetUserID = 0
		session.TargetGroupID = reciveMsg.GroupID

		var iSession services.ISession
		iSession = &session
		services.OnRecivedMessage(iSession)

	}
}
