package nasa

import (
	"Asane/internal/system"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type client struct {
	APIKey string
}

var apiURL = url.URL{
	Scheme: "https",
	Host:   "api.nasa.gov",
}

// httpClient 单例
var httpClient = &client{
	APIKey: system.Config.NasaConfig.APIKey,
}

// APOD NASA每日最佳图片
func (c *client) APOD(date string) (APODResponseObject, error) {
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
