package qq

import (
	"net/http"
	"strings"

	"github.com/NagaseYami/asane/system"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

var Server = server{
	bots: make(map[string]*bot),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type server struct {
	address string
	token   string
	bots    map[string]*bot
}

func (s *server) Run() {
	s.address = system.Config.Get("qq_config.address").String()
	s.token = system.Config.Get("qq_config.token").String()

	http.HandleFunc("/", s.anyMessage)
	log.Infof("go-cqhttp专用反向Websocket服务器已启动: %v", s.address)
	log.Fatal(http.ListenAndServe(s.address, nil))
}

func (s *server) anyMessage(w http.ResponseWriter, r *http.Request) {
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
	bot := &bot{
		id:   botID,
		conn: c,
	}
	s.bots[botID] = bot
	go s.listenEvent(bot)
}

func (s *server) authorization(r *http.Request) bool {
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

func (s *server) listenEvent(b *bot) {
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		system.HandleError(err)
	}(b.conn)

	defer delete(s.bots, b.id)

	for {
		t, payload, err := b.conn.ReadMessage()
		if err != nil {
			log.Warnf("监听反向WS API时出现错误: %v", err)
			break
		}
		if t == websocket.TextMessage {
			b.universalRouter(gjson.ParseBytes(payload))
		}
	}
}
