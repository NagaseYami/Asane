package qq

import (
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

type websocketServer struct {
	token string
	conn  *websocket.Conn
}

// WebSocketServer Singleton
var WebSocketServer = &websocketServer{
	token: os.Getenv("WEBSOCKET_SERVER_TOKEN"),
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
	s.conn = c
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
