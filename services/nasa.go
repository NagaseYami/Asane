package services

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/NagaseYami/asane/system"
	log "github.com/sirupsen/logrus"
)

var apodApiURL = url.URL{
	Scheme: "https",
	Host:   "api.nasa.gov",
	Path:   "planetary/apod",
}

func APOD(params []string) (gjson.Result, error) {
	date := ""

	if len(params) > 0 {
		date = params[0]
		match, err := regexp.MatchString("[0-9]{4}-[0-9]{2}-[0-9]{2}", date)
		system.HandleError(err)

		if match {
			date = params[0]
		} else {
			return gjson.Result{}, errors.New("日期指定格式错误。正确格式：YYYY-MM-DD")
		}
	}

	apodURL := apodApiURL
	query := url.Values{}
	if date != "" {
		query.Add("date", date)
	}
	query.Add("api_key", system.Config.Get("nasa_config.api_key").String())
	apodURL.RawQuery = query.Encode()

	resp, err := http.Get(apodURL.String())
	if err != nil {
		log.Error(err)
		return gjson.Result{}, errors.New("发生未知错误，请检查日志")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return gjson.Result{}, fmt.Errorf("NASA没有%s这一天的APOD", date)
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, resp.Body)

	return gjson.ParseBytes(buf.Bytes()), nil
}
