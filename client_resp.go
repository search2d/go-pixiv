package pixiv

type APIErrorBody struct {
	Error APIError `json:"error"`
}

type APIError struct {
	UserMessage        string                 `json:"user_message"`
	Message            string                 `json:"message"`
	Reason             string                 `json:"reason"`
	UserMessageDetails map[string]interface{} `json:"user_message_details"`
}

type GetIllustRanking struct {
	Illusts []GetIllustRankingIllust `json:"illusts"`
	NextURL string                   `json:"next_url"`
}

type GetIllustRankingIllust struct {
	ID             int                              `json:"id"`
	Title          string                           `json:"title"`
	Type           string                           `json:"type"`
	ImageURLs      map[string]string                `json:"image_urls"`
	Caption        string                           `json:"caption"`
	Restrict       int                              `json:"restrict"`
	User           GetIllustRankingIllustUser       `json:"user"`
	Tags           []GetIllustRankingIllustTag      `json:"tags"`
	Tools          []string                         `json:"tools"`
	CreateDate     string                           `json:"create_date"`
	PageCount      int                              `json:"page_count"`
	Width          int                              `json:"width"`
	Height         int                              `json:"height"`
	SanityLevel    int                              `json:"sanity_level"`
	Series         GetIllustRankingIllustSeries     `json:"series"`
	MetaSinglePage map[string]string                `json:"meta_single_page"`
	MetaPages      []GetIllustRankingIllustMetaPage `json:"meta_pages"`
	TotalView      int                              `json:"total_view"`
	TotalBookmarks int                              `json:"total_bookmarks"`
	IsBookmarked   bool                             `json:"is_bookmarked"`
	Visible        bool                             `json:"visible"`
	IsMuted        bool                             `json:"is_muted"`
}

type GetIllustRankingIllustUser struct {
	ID               int               `json:"id"`
	Name             string            `json:"name"`
	Account          string            `json:"account"`
	ProfileImageURLs map[string]string `json:"profile_image_urls"`
	IsFollowed       bool              `json:"is_followed"`
}

type GetIllustRankingIllustTag struct {
	Name string `json:"name"`
}

type GetIllustRankingIllustSeries struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type GetIllustRankingIllustMetaPage struct {
	ImageURLs map[string]string `json:"image_urls"`
}

type GetIllustDetail struct {
	Illust GetIllustDetailIllust `json:"illust"`
}

type GetIllustDetailIllust struct {
	ID             int                             `json:"id"`
	Title          string                          `json:"title"`
	Type           string                          `json:"type"`
	ImageURLs      map[string]string               `json:"image_urls"`
	Caption        string                          `json:"caption"`
	Restrict       int                             `json:"restrict"`
	User           GetIllustDetailIllustUser       `json:"user"`
	Tags           []GetIllustDetailIllustTag      `json:"tags"`
	Tools          []string                        `json:"tools"`
	CreateDate     string                          `json:"create_date"`
	PageCount      int                             `json:"page_count"`
	Width          int                             `json:"width"`
	Height         int                             `json:"height"`
	SanityLevel    int                             `json:"sanity_level"`
	Series         GetIllustDetailIllustSeries     `json:"series"`
	MetaSinglePage map[string]string               `json:"meta_single_page"`
	MetaPages      []GetIllustDetailIllustMetaPage `json:"meta_pages"`
	TotalView      int                             `json:"total_view"`
	TotalBookmarks int                             `json:"total_bookmarks"`
	IsBookmarked   bool                            `json:"is_bookmarked"`
	Visible        bool                            `json:"visible"`
	IsMuted        bool                            `json:"is_muted"`
	TotalComments  int                             `json:"total_comments"`
}

type GetIllustDetailIllustUser struct {
	ID               int               `json:"id"`
	Name             string            `json:"name"`
	Account          string            `json:"alice810"`
	ProfileImageURLs map[string]string `json:"profile_image_urls"`
	IsFollowed       bool              `json:"is_followed"`
}

type GetIllustDetailIllustTag struct {
	Name string `json:"name"`
}

type GetIllustDetailIllustSeries struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type GetIllustDetailIllustMetaPage struct {
	ImageURLs map[string]string `json:"image_urls"`
}
