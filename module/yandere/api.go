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

// PostsListResponseObject PostList Response
type PostsListResponseObject struct {
	ID                  int           `json:"id"`
	Tags                string        `json:"tags"`
	CreatedAt           int           `json:"created_at"`
	UpdatedAt           int           `json:"updated_at"`
	CreatorID           int           `json:"creator_id"`
	ApproverID          interface{}   `json:"approver_id"`
	Author              string        `json:"author"`
	Change              int           `json:"change"`
	Source              string        `json:"source"`
	Score               int           `json:"score"`
	Md5                 string        `json:"md5"`
	FileSize            int           `json:"file_size"`
	FileExt             string        `json:"file_ext"`
	FileURL             string        `json:"file_url"`
	IsShownInIndex      bool          `json:"is_shown_in_index"`
	PreviewURL          string        `json:"preview_url"`
	PreviewWidth        int           `json:"preview_width"`
	PreviewHeight       int           `json:"preview_height"`
	ActualPreviewWidth  int           `json:"actual_preview_width"`
	ActualPreviewHeight int           `json:"actual_preview_height"`
	SampleURL           string        `json:"sample_url"`
	SampleWidth         int           `json:"sample_width"`
	SampleHeight        int           `json:"sample_height"`
	SampleFileSize      int           `json:"sample_file_size"`
	JpegURL             string        `json:"jpeg_url"`
	JpegWidth           int           `json:"jpeg_width"`
	JpegHeight          int           `json:"jpeg_height"`
	JpegFileSize        int           `json:"jpeg_file_size"`
	Rating              string        `json:"rating"`
	IsRatingLocked      bool          `json:"is_rating_locked"`
	HasChildren         bool          `json:"has_children"`
	ParentID            interface{}   `json:"parent_id"`
	Status              string        `json:"status"`
	IsPending           bool          `json:"is_pending"`
	Width               int           `json:"width"`
	Height              int           `json:"height"`
	IsHeld              bool          `json:"is_held"`
	FramesPendingString string        `json:"frames_pending_string"`
	FramesPending       []interface{} `json:"frames_pending"`
	FramesString        string        `json:"frames_string"`
	Frames              []interface{} `json:"frames"`
	IsNoteLocked        bool          `json:"is_note_locked"`
	LastNotedAt         int           `json:"last_noted_at"`
	LastCommentedAt     int           `json:"last_commented_at"`
}

// PostsListRequestQueryObject PostsList Request Query
type PostsListRequestQueryObject struct {
	Limit int
	Page  int
	Tags  string
}

// URL Get API URL
func (api *PostsListRequestQueryObject) URL() *url.URL {
	result := apiURL
	result.Path = path.Join(result.Path, "post.json")
	query := result.Query()

	if api.Limit != 0 {
		query.Add("limit", strconv.Itoa(api.Limit))
	}
	if api.Page != 0 {
		query.Add("page", strconv.Itoa(api.Page))
	}
	if api.Tags != "" {
		query.Add("tags", api.Tags)
	}

	result.RawQuery = query.Encode()

	return &result
}
