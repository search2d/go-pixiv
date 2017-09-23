package pixiv

type TokenErrorBody struct {
	HasError bool                  `json:"has_error"`
	Errors   map[string]TokenError `json:"errors"`
}

type TokenError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Token struct {
	Response TokenResponse `json:"response"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	ExpiresIn    int       `json:"expires_in"`
	TokenType    string    `json:"token_type"`
	Scope        string    `json:"scope"`
	RefreshToken string    `json:"refresh_token"`
	User         TokenUser `json:"user"`
}

type TokenUser struct {
	ProfileImageURLs map[string]string `json:"profile_image_urls"`
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Account          string            `json:"account"`
	MailAddress      string            `json:"mail_address"`
	IsPremium        bool              `json:"is_premium"`
	XRestrict        int               `json:"x_restrict"`
	IsMailAuthorized bool              `json:"is_mail_authorized"`
}
