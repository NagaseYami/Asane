package yandere

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type client struct{}

var yandereURL = url.URL{
	Scheme: "https",
	Host:   "yande.re",
}

// Client YandereClient单例
var Client = &client{}

func (c *client) GetRandomExplicitPost() YanderePostsListResponseObject {
	log.Info("调用yandere的随机色图api")
	defer log.Info("成功获取并返回随机色图信息")
	api := &YanderePostsListApi{
		Limit: 1,
		Tags:  "rating:explicit order:random score:>30",
	}
	resp, err := http.Get(api.GetURL().String())
	if err != nil {
		log.Warnln(err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, resp.Body)

	yanderePostsListResponse := &[]YanderePostsListResponseObject{}
	err = json.Unmarshal(buf.Bytes(), yanderePostsListResponse)
	if err != nil {
		log.Panic(err)
	}
	return (*yanderePostsListResponse)[0]
}
