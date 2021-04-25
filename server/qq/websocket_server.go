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
	token          string
	messageHandler func(*websocket.Conn, IReciveMessageObject)
}

// WebSocketServer Singleton
var WebSocketServer = &websocketServer{
	token:          os.Getenv("WEBSOCKET_SERVER_TOKEN"),
	messageHandler: func(*websocket.Conn, IReciveMessageObject) {},
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

func (s *websocketServer) HandleMessage(handler func(*websocket.Conn, IReciveMessageObject)) {
	s.messageHandler = handler
}

func (s *websocketServer) WriteTextMessage(conn *websocket.Conn, b []byte) {
	log.Info("发送bytes到CqHttp")
	conn.WriteMessage(websocket.TextMessage, b)
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
			s.universalRouter(conn, gjson.ParseBytes(payload))
		}
	}
}

func (s *websocketServer) universalRouter(conn *websocket.Conn, result gjson.Result) {
	switch result.Get("post_type").Str {
	case "message":
		log.Trace(result.String())
		s.messageRouter(conn, result)

	case "notice":

	case "request":

	case "meta_event":
		s.metaEventRouter(conn, result)

	}
}

func (s *websocketServer) metaEventRouter(conn *websocket.Conn, result gjson.Result) {
	switch result.Get("meta_event_type").Str {
	case "lifecycle":

	case "heartbeat":

	}
}

func (s *websocketServer) messageRouter(conn *websocket.Conn, result gjson.Result) {
	switch result.Get("message_type").Str {
	case "private":
		reciveMsg := &RecivePrivateMessageObject{}
		json.Unmarshal([]byte(result.Raw), reciveMsg)
		var iMsg IReciveMessageObject
		iMsg = reciveMsg
		s.messageHandler(conn, iMsg)

	case "group":
		reciveMsg := &ReciveGroupMessageObject{}
		json.Unmarshal([]byte(result.Raw), reciveMsg)
		var iMsg IReciveMessageObject
		iMsg = reciveMsg
		s.messageHandler(conn, iMsg)

	}
}
