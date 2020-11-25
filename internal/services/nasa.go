package services

import (
	"Asane/internal/api/nasa"
	"fmt"
	"regexp"
)

func nasaAPOD(params []string) string {
	date := ""
	if len(params) > 0 {
		date = params[0]
		if match, _ := regexp.Match("[0-9]{4}-[0-9]{2}-[0-9]{2}", []byte(date)); match {
			date = params[0]
		} else {
			return " 日期指定格式错误。正确格式：YYYY-MM-DD"
		}
	}

	resp, err := nasa.Client.APOD(date)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("%s\n%s\n[CQ:image,file=%s]", resp.Title, resp.Explanation, resp.URL)
}
