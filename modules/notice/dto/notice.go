package dto

type NoticeChannel struct {
	Uuid  string `json:"uuid"`
	Title string `json:"title"`
	Type  int    `json:"type"` //1邮箱 2webhook
}
