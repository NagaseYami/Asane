package services

import (
	"Asane/internal/api/nasa"
	"fmt"
)

func nasaAPOD(params []string) string {
	date := ""
	if len(params) > 0 {
		date = params[0]
	}

	resp,err := nasa.Client.APOD(date)
	if err!=nil{
		return err.Error()
	}

	return fmt.Sprintf("%s\n%s\n[CQ:image,file=%s]", resp.Title, resp.Explanation, resp.URL)
}
