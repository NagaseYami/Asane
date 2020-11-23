package yandere

import (
	"net/url"
	"path"
	"strconv"
)

type YandereTagsListResponseObject struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Count     int    `json:"count"`
	Type      int    `json:"type"`
	Ambiguous bool   `json:"ambiguous"`
}

type YandereTagsListApi struct {
	Limit int
	Order string
	Name  string
}

func (api *YandereTagsListApi) GetURL() *url.URL {
	result := yandereURL
	result.Path = path.Join(result.Path, "tag.json")
	query := result.Query()

	if api.Limit != 0 {
		query.Add("limit", strconv.Itoa(api.Limit))
	}
	if api.Order != "" {
		query.Add("order", api.Order)
	}
	if api.Name != "" {
		query.Add("name", api.Name)
	}

	result.RawQuery = query.Encode()

	return &result
}
