package sankaku

import (
	"net/url"
	"path"
	"strconv"
)

type TagListResponseObject struct {
	ID          int           `json:"id"`
	NameEn      string        `json:"name_en"`
	NameJa      string        `json:"name_ja"`
	Type        int           `json:"type"`
	Count       int           `json:"count"`
	PostCount   int           `json:"post_count"`
	PoolCount   int           `json:"pool_count"`
	Locale      string        `json:"locale"`
	Rating      string        `json:"rating"`
	Name        string        `json:"name"`
	RelatedTags []interface{} `json:"related_tags"`
	ChildTags   []interface{} `json:"child_tags"`
	ParentTags  []interface{} `json:"parent_tags"`
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
	result.Path = path.Join(result.Path, "tags")
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
