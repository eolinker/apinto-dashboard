package dto

type WebhookInput struct {
	Uuid          string            `json:"uuid"`
	Title         string            `json:"title"`
	Desc          string            `json:"desc"`
	Url           string            `json:"url"`
	Method        string            `json:"method"`
	ContentType   string            `json:"content_type"` //
	NoticeType    string            `json:"notice_type"`  //单次：single  多次：many
	UserSeparator string            `json:"user_separator"`
	Header        map[string]string `json:"header"`
	Template      string            `json:"template"`
}

type WebhooksOutput struct {
	Uuid        string `json:"uuid"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	Method      string `json:"method"`
	ContentType string `json:"content_type"`
	Operator    string `json:"operator"`
	UpdateTime  string `json:"update_time"`
	CreateTime  string `json:"create_time"`
	IsDelete    bool   `json:"is_delete"`
}

type WebhookOutput struct {
	Uuid          string            `json:"uuid"`
	Title         string            `json:"title"`
	Desc          string            `json:"desc"`
	Url           string            `json:"url"`
	Method        string            `json:"method"`
	ContentType   string            `json:"content_type"` //
	NoticeType    string            `json:"notice_type"`  //单次：single  多次：many
	UserSeparator string            `json:"user_separator"`
	Header        map[string]string `json:"header"`
	Template      string            `json:"template"`
}
