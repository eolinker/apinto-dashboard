package user_model

import (
	user_entry "github.com/eolinker/apinto-dashboard/modules/user/user-entry"
	"time"
)

type UserInfo struct {
	Id            int        `json:"id"`                       //用户ID
	Sex           int        `json:"sex,omitempty"`            //0未知 1男2女
	UserName      string     `json:"user_name,omitempty"`      //账号
	NoticeUserId  string     `json:"notice_user_id,omitempty"` //通知用户ID
	NickName      string     `json:"nick_name,omitempty"`      //用户昵称
	Email         string     `json:"email,omitempty"`          //邮箱
	Phone         string     `json:"phone,omitempty"`          //手机号
	Avatar        string     `json:"avatar,omitempty"`         //头像
	LastLoginTime *time.Time `json:"login_time,omitempty"`     //最近登录时间
}
type UserBase struct {
	Sex          int    `json:"sex,omitempty"`            //0未知 1男2女
	UserName     string `json:"user_name,omitempty"`      //账号
	NoticeUserId string `json:"notice_user_id,omitempty"` //通知用户ID
	NickName     string `json:"nick_name,omitempty"`      //用户昵称
	Email        string `json:"email,omitempty"`          //邮箱
	Phone        string `json:"phone,omitempty"`          //手机号
	Avatar       string `json:"avatar,omitempty"`         //头像
}

func CreateUserInfo(info *user_entry.UserInfo) *UserInfo {
	return &UserInfo{
		Id:            info.Id,
		Sex:           info.Sex,
		UserName:      info.UserName,
		NoticeUserId:  info.NoticeUserId,
		NickName:      info.NickName,
		Email:         info.Email,
		Phone:         info.Phone,
		Avatar:        info.Avatar,
		LastLoginTime: info.LastLoginTime,
	}
}
