package nasa

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/NagaseYami/asane/system"
	log "github.com/sirupsen/logrus"
)

func APOD(params []string) (APODResponseObject, error) {
	date := ""

	if len(params) > 0 {
		date = params[0]
		match, err := regexp.MatchString("[0-9]{4}-[0-9]{2}-[0-9]{2}", date)
		if err != nil {
			log.Panic(err)
		}

		if match {
			date = params[0]
		} else {
			return APODResponseObject{}, errors.New("日期指定格式错误。正确格式：YYYY-MM-DD")
		}
	}

	api := &APODRequestQueryObject{
		Date:   date,
		HD:     false,
		APIKey: system.Config.NasaConfig.APIKey,
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
