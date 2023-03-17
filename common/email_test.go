package common

import (
	"gopkg.in/gomail.v2"
	"testing"
)

func TestNewSmtp(t *testing.T) {
	type args struct {
		host     string
		port     int
		protocol string
		account  string
		password string
		email    string
	}
	tests := []struct {
		name string
		args args
		want *Smtp
	}{
		{
			name: "",
			args: args{
				host:     "smtp.qq.com",
				port:     587,
				protocol: "ssl",
				account:  "1324204490@qq.com",
				password: "zzeqxoubzzoababg",
				email:    "1324204490@qq.com",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			smtp := NewSmtp(tt.args.host, tt.args.port, tt.args.protocol, tt.args.account, tt.args.password, tt.args.email)
			sendTo := make([]string, 0)
			sendTo = append(sendTo, "zhangzeyi4490@foxmail.com", "1324204490@qq.com")
			if err := smtp.SendTo(sendTo, "告警策略1", "请求失败状态码数/5分钟 > 5次\\n        <span style=\\\"color: #f04864\\\">实际值：10次</span>；且 请求成功率/5分钟\\n        同比波动 5% <span style=\\\"color: #f04864\\\">实际值：10%</span>;\\n        <br />或<br />\\n        请求失败状态码数/5分钟 同比波动 5次"); err != nil {
				t.Errorf("sendTo() , want %s", err.Error())
			}

		})
	}
}

func TestSmtp_SendTo(t *testing.T) {
	type fields struct {
		email  string
		dialer *gomail.Dialer
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
			s := &Smtp{
				email:  tt.fields.email,
				dialer: tt.fields.dialer,
			}
			if err := s.SendTo(tt.args.sends, tt.args.title, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("SendTo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
