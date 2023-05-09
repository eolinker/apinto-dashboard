package common

import (
	"net/http"
	"reflect"
	"testing"
)

func TestNewWebhook(t *testing.T) {
	type args struct {
		uri           string
		method        string
		contentType   string
		noticeType    string
		userSeparator string
		header        map[string]string
		template      string
	}
	tests := []struct {
		name string
		args args
		want *Webhook
	}{
		{
			name: "",
			args: args{
				uri:           "https://open.feishu.cn/open-apis/bot/v2/hook/a7ed8efa-88ac-4721-af0e-c00d02172312",
				method:        http.MethodPost,
				contentType:   "JSON",
				noticeType:    SingleNotice,
				userSeparator: ",",
				header:        map[string]string{"test1": "test2"},
				template: `{
    "msg_type": "interactive",
    "card": {
        "elements": [{
                "tag": "div",
                "text": {
                        "content": "{msg}",
                        "tag": "lark_md"
                }
        }],
        "header": {
                "title": {
                        "content": "{title}",
                        "tag": "plain_text"
                }
        }
    }
}`,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			webhook := NewWebhook(tt.args.uri, tt.args.method, tt.args.contentType, tt.args.noticeType, tt.args.userSeparator, tt.args.header, tt.args.template)
			sendTo := make([]string, 0)
			sendTo = append(sendTo, "1324204490", "18273886975")

			if err := webhook.SendTo(sendTo, "API告警策略", "你好你好你好"); err != nil {
				t.Errorf("SendTo() , want %v err=%s", tt.want, err.Error())
			}
		})
	}
}

func TestWebhook_SendTo(t *testing.T) {
	type fields struct {
		Url           string
		Method        string
		ContentType   string
		NoticeType    string
		UserSeparator string
		Header        map[string]string
		Template      string
	}
	type args struct {
		sends []string
		title string
		msg   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Webhook{
				Url:           tt.fields.Url,
				Method:        tt.fields.Method,
				ContentType:   tt.fields.ContentType,
				NoticeType:    tt.fields.NoticeType,
				UserSeparator: tt.fields.UserSeparator,
				Header:        tt.fields.Header,
				Template:      tt.fields.Template,
			}
			if err := w.SendTo(tt.args.sends, tt.args.title, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("SendTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWebhook_genUsers(t *testing.T) {
	type fields struct {
		Url           string
		Method        string
		ContentType   string
		NoticeType    string
		UserSeparator string
		Header        map[string]string
		Template      string
	}
	type args struct {
		users []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Webhook{
				Url:           tt.fields.Url,
				Method:        tt.fields.Method,
				ContentType:   tt.fields.ContentType,
				NoticeType:    tt.fields.NoticeType,
				UserSeparator: tt.fields.UserSeparator,
				Header:        tt.fields.Header,
				Template:      tt.fields.Template,
			}
			if got := w.genUsers(tt.args.users); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("genUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_replaceTemplate(t *testing.T) {
	type args struct {
		org   string
		title string
		msg   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceTemplate(tt.args.org, tt.args.title, tt.args.msg); got != tt.want {
				t.Errorf("replaceTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_replaceUsers(t *testing.T) {
	type args struct {
		org   string
		users string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceUsers(tt.args.org, tt.args.users); got != tt.want {
				t.Errorf("replaceUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}
