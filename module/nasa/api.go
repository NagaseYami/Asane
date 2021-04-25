package nasa

import (
	"net/url"
	"path"
	"strconv"
)

// APODResponseObject Nasa APOD api response
type APODResponseObject struct {
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	Hdurl          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	URL            string `json:"url"`
}

// APODRequestQueryObject Nasa APOD api request query
type APODRequestQueryObject struct {
	Date   string
	HD     bool
	APIKey string
}

// URL Get API URL
func (api *APODRequestQueryObject) URL() *url.URL {
	result := apiURL
	result.Path = path.Join(result.Path, "planetary/apod")
	query := result.Query()

	if api.Date != "" {
		query.Add("date", api.Date)
	}
	query.Add("hd", strconv.FormatBool(api.HD))
	query.Add("api_key", api.APIKey)

	result.RawQuery = query.Encode()

	return &result
}
