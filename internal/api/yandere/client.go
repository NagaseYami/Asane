package yandere

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

var yandereURL = url.URL{
	Scheme: "https",
	Host:   "yande.re",
}

// Client YandereClient单例
var Client = &client{}

func (c *client) SearchTags(tag string) ([]YandereTagsListResponseObject, error) {
	api := &YandereTagsListApi{
		Limit: 10,
		Name:  tag,
		Order: "count",
	}
	resp, err := http.Get(api.GetURL().String())
	if err != nil {
		log.Warnln()
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, resp.Body)

	result := &[]YandereTagsListResponseObject{}
	err = json.Unmarshal(buf.Bytes(), result)
	if err != nil {
		log.Fatal(err)
	}

	if len(*result) == 0 {
		return []YandereTagsListResponseObject{}, errors.New("搜索结果为0")
	}

	return (*result), nil
}

func (c *client) GetRandomExplicitPost(tags []string) (YanderePostsListResponseObject, error) {
	if tags == nil {
		tags = []string{}
	}

	tags = append(tags, "rating:explicit")
	tags = append(tags, "order:random")

	api := &YanderePostsListApi{
		Limit: 1,
		Tags:  strings.Join(tags, " "),
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
		log.Fatal(err)
	}
	if len(*yanderePostsListResponse) == 0 {
		return YanderePostsListResponseObject{}, errors.New("搜索结果为0")
	}
	return (*yanderePostsListResponse)[0], nil
}
