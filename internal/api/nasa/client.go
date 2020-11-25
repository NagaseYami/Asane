package nasa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
)

type client struct {
	APIKey string
}

var apiURL = url.URL{
	Scheme: "https",
	Host:   "api.nasa.gov",
}

// Client NASAClient单例
var Client = &client{
	APIKey: os.Getenv("NASA_API_KEY"),
}

// APOD NASA每日最佳图片
func (c *client) APOD(date string) (APODResponseObject, error) {
	if c.APIKey == "" {
		return APODResponseObject{}, errors.New("缺少环境变量NASA_API_KEY，该功能无法使用")
	}

	api := &APODRequestQueryObject{
		Date:   date,
		HD:     false,
		APIKey: c.APIKey,
	}

	resp, err := http.Get(api.URL().String())
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return APODResponseObject{}, fmt.Errorf("NASA没有%s这一天的APOD", date)
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, resp.Body)

	result := &APODResponseObject{}
	err = json.Unmarshal(buf.Bytes(), result)

	if err != nil {
		log.Error(err)
	}

	return *result, nil
}
