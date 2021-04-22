package yandere

import (
	"net/url"
	"path"
	"strconv"
)

// TagListResponseObject TagList Response
type TagListResponseObject struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Count     int    `json:"count"`
	Type      int    `json:"type"`
	Ambiguous bool   `json:"ambiguous"`
}

// TagListRequestQueryObject TagList Request Query
type TagListRequestQueryObject struct {
	Limit int
	Order string
	Name  string
}

// URL Get API URL
func (api *TagListRequestQueryObject) URL() *url.URL {
	result := apiURL
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
