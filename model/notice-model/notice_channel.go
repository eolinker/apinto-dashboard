package notice_model

import (
	"github.com/eolinker/apinto-dashboard/entry/notice-entry"
	"time"
)

type NoticeChannel struct {
	Id         int
	Name       string
	Title      string
	Type       int    //1webhook 2mail
	Config     string // NoticeChannelEmail NoticeChannelWebHook 转成json
	IsDelete   bool
	Operator   string
	CreateTime time.Time
	UpdateTime time.Time
}

type NoticeChannelVersion notice_entry.NoticeChannelVersion

type NoticeChannelWebhook struct {
	Desc          string            `json:"desc"`
	Url           string            `json:"url"`
	Method        string            `json:"method"`
	ContentType   string            `json:"content_type"`
	NoticeType    string            `json:"notice_type"`
	UserSeparator string            `json:"user_separator"`
	Header        map[string]string `json:"header"`
	Template      string            `json:"template"`
}

type NoticeChannelEmail struct {
	SmtpUrl  string `json:"smtp_url"`
	SmtpPort int    `json:"smtp_port"`
	Protocol string `json:"protocol"`
	Account  string `json:"account"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
