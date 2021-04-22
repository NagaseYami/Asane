package sankaku

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

type client struct{}

var apiURL = url.URL{
	Scheme: "https",
	Host:   "capi-v2.sankakucomplex.com",
}

// httpClient 单例
var httpClient = &client{}

func (c *client) SearchTags(tag string) ([]TagListResponseObject, error) {
	api := &TagListRequestQueryObject{
		Limit: 10,
		Name:  tag,
		Order: "count",
	}
	resp, err := http.Get(api.URL().String())
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, resp.Body)

	result := &[]TagListResponseObject{}
	err = json.Unmarshal(buf.Bytes(), result)
	if err != nil {
		log.Error(err)
	}

	if len(*result) == 0 {
		return []TagListResponseObject{}, errors.New("搜索结果为0")
	}

	return (*result), nil
}

func (c *client) RandomExplicitPost(tags []string) (PostsListResponseObject, error) {
	if tags == nil {
		tags = []string{}
	}

	tags = append(tags, "rating:explicit")
	tags = append(tags, "order:random")

	api := &PostsListRequestQueryObject{
		Limit: 1,
		Tags:  strings.Join(tags, " "),
	}
	resp, err := http.Get(api.URL().String())
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, resp.Body)

	postsListResponse := &[]PostsListResponseObject{}
	err = json.Unmarshal(buf.Bytes(), postsListResponse)
	if err != nil {
		log.Error(err)
	}
	if len(*postsListResponse) == 0 {
		return PostsListResponseObject{}, errors.New("搜索结果为0\n本功能仅支持精准tag搜索\n请确认您输入的tag拼写正确\n您可以使用asane tag命令来搜索tag")
	}
	return (*postsListResponse)[0], nil
}
