package common

import (
	"github.com/eolinker/eosc/log"
	"io"
	"net/http"
	"strings"
)

var client = http.DefaultClient

const (
	SingleNotice = "single"
)

type Webhook struct {
	Url           string
	Method        string
	ContentType   string
	NoticeType    string
	UserSeparator string
	Header        map[string]string
	Template      string
}

func NewWebhook(uri, method string, contentType string, noticeType, userSeparator string, header map[string]string, template string) *Webhook {
	switch contentType {
	case "JSON":
		{
			header["content-type"] = "application/json"
			if template == "" {
				template = `{"users":"{user}","title":"{title}","msg":"{msg}"}`
			}
		}
	case "form-data":
		{
			header["content-type"] = "application/x-www-form-urlencoded"
			if template == "" {
				template = "users={users}&title={title}&msg={msg}"
			}
		}
	}
	return &Webhook{
		Url:           uri,
		Method:        method,
		ContentType:   contentType,
		NoticeType:    noticeType,
		UserSeparator: userSeparator,
		Header:        header,
		Template:      template,
	}
}

func send(method string, uri string, body string, headers map[string]string) error {
	request, err := http.NewRequest(method, uri, strings.NewReader(body))
	if err != nil {
		return err
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}
	bytes, _ := io.ReadAll(resp.Body)
	log.Errorf("Webhook-sendTo uri=%s requestHeader=%v reqBody=%s respCode=%v respBody=%s", uri, headers, body, resp.StatusCode, string(bytes))
	return nil
}

func (w *Webhook) SendTo(sends []string, title, msg string) error {
	body := replaceTemplate(w.Template, title, msg)
	uri := replaceTemplate(w.Url, title, msg)
	if len(sends) < 1 {
		return send(w.Method, uri, body, w.Header)
	}

	users := w.genUsers(sends)
	for _, user := range users {
		err := send(w.Method, replaceUsers(uri, user), replaceUsers(body, user), w.Header)
		if err != nil {
			log.Errorf("send request error: %v,user is %s", err, user)
			return err
		}
	}
	return nil
}

func (w *Webhook) genUsers(users []string) []string {
	if w.NoticeType == SingleNotice {
		return []string{strings.Join(users, w.UserSeparator)}
	}
	return users
}

func replaceUsers(org, users string) string {
	return strings.Replace(org, "{users}", users, -1)
}

func replaceTemplate(org, title, msg string) string {
	org = strings.Replace(org, "{title}", title, -1)
	org = strings.Replace(org, "{msg}", msg, -1)
	return org
}
