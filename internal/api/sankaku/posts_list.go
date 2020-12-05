package sankaku

import (
	"net/url"
	"path"
	"strconv"
)

type PostsListResponseObject struct {
	ID     int    `json:"id"`
	Rating string `json:"rating"`
	Status string `json:"status"`
	Author struct {
		ID           int    `json:"id"`
		Name         string `json:"name"`
		Avatar       string `json:"avatar"`
		AvatarRating string `json:"avatar_rating"`
	} `json:"author"`
	SampleURL     string `json:"sample_url"`
	SampleWidth   int    `json:"sample_width"`
	SampleHeight  int    `json:"sample_height"`
	PreviewURL    string `json:"preview_url"`
	PreviewWidth  int    `json:"preview_width"`
	PreviewHeight int    `json:"preview_height"`
	FileURL       string `json:"file_url"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	FileSize      int    `json:"file_size"`
	FileType      string `json:"file_type"`
	CreatedAt     struct {
		JSONClass string `json:"json_class"`
		S         int    `json:"s"`
		N         int    `json:"n"`
	} `json:"created_at"`
	HasChildren      bool        `json:"has_children"`
	HasComments      bool        `json:"has_comments"`
	HasNotes         bool        `json:"has_notes"`
	IsFavorited      bool        `json:"is_favorited"`
	UserVote         interface{} `json:"user_vote"`
	Md5              string      `json:"md5"`
	ParentID         int         `json:"parent_id"`
	Change           int         `json:"change"`
	FavCount         int         `json:"fav_count"`
	RecommendedPosts int         `json:"recommended_posts"`
	RecommendedScore int         `json:"recommended_score"`
	VoteCount        int         `json:"vote_count"`
	TotalScore       int         `json:"total_score"`
	CommentCount     interface{} `json:"comment_count"`
	Source           string      `json:"source"`
	InVisiblePool    bool        `json:"in_visible_pool"`
	IsPremium        bool        `json:"is_premium"`
	Sequence         interface{} `json:"sequence"`
	Tags             []struct {
		ID        int         `json:"id"`
		NameEn    string      `json:"name_en"`
		NameJa    string      `json:"name_ja"`
		Type      int         `json:"type"`
		Count     int         `json:"count"`
		PostCount int         `json:"post_count"`
		PoolCount int         `json:"pool_count"`
		Locale    string      `json:"locale"`
		Rating    interface{} `json:"rating"`
		Name      string      `json:"name"`
	} `json:"tags"`
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
	result.Path = path.Join(result.Path, "posts")
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
