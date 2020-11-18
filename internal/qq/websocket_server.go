package qq

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type websocketServer struct {
	token string
	conn  *websocket.Conn

	privateMessageHandler map[string]func(gjson.Result)
	groupMessageHandler   map[string]func(gjson.Result)
}

// WebSocketServer Singleton
var WebSocketServer = &websocketServer{
	token:                 os.Getenv("WEBSOCKET_SERVER_TOKEN"),
	conn:                  nil,
	privateMessageHandler: map[string]func(gjson.Result){},
	groupMessageHandler:   map[string]func(gjson.Result){},
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *websocketServer) Run(addr, authToken string) {
	s.token = authToken
	http.HandleFunc("/", s.any)
	log.Infof("CQ Websocket 服务器已启动: %v", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func (s *websocketServer) HandlePrivateMessageFunc(keyword string, handler func(gjson.Result)) {
	s.privateMessageHandler[keyword] = handler
}

func (s *websocketServer) HandleGroupMessageFunc(keyword string, handler func(gjson.Result)) {
	s.groupMessageHandler[keyword] = handler
}

func (s *websocketServer) SendJSON(message interface{}) {
	log.Info("发送Json至CqHttp")
	b, err := json.Marshal(message)
	if err != nil {
		log.Panic(err)
	}
	s.conn.WriteMessage(websocket.TextMessage, b)
}

func (s *websocketServer) any(w http.ResponseWriter, r *http.Request) {
	if !s.authorization(r) {
		log.Warnf("已拒绝 %v 的 Websocket 请求: Token鉴权失败", r.RemoteAddr)
		w.WriteHeader(401)
	}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warnf("处理 Websocket 请求时出现错误: %v", err)
		return
	}
	log.Infof("接受 Websocket 连接: %v (/)", r.RemoteAddr)
	// FIXME
	s.conn = c
	go s.listenEvent(c)
}

func (s *websocketServer) authorization(r *http.Request) bool {
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

func (s *websocketServer) listenEvent(conn *websocket.Conn) {
	defer conn.Close()
	for {
		t, payload, err := conn.ReadMessage()
		if err != nil {
			log.Warnf("监听反向WS API时出现错误: %v", err)
			break
		}
		if t == websocket.TextMessage {
			go s.universalHandler(gjson.ParseBytes(payload))
		}
	}
}

func (s *websocketServer) universalHandler(json gjson.Result) {
	switch json.Get("post_type").Str {
	case "message":
		log.Debug(json.String())
		s.messageHandler(json)
		break
	case "notice":
		break
	case "request":
		break
	case "meta_event":
		s.metaEventHandler(json)
		break
	}
}

func (s *websocketServer) metaEventHandler(json gjson.Result) {
	switch json.Get("meta_event_type").Str {
	case "lifecycle":
		break
	case "heartbeat":
		break
	}
}

func (s *websocketServer) messageHandler(json gjson.Result) {
	switch json.Get("message_type").Str {
	case "private":
		fn, exist := s.privateMessageHandler[json.Get("raw_message").Str]
		if exist {
			fn(json)
		}
		break
	case "group":
		fn, exist := s.groupMessageHandler[json.Get("raw_message").Str]
		if exist {
			fn(json)
		}
		break
	}
}
