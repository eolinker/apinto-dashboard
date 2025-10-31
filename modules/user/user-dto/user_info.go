package user_dto

type UserEnum struct {
	Sex          int    `json:"sex,omitempty"`          //0未知 1男2女
	UserName     string `json:"userName,omitempty"`     //账号
	NoticeUserId string `json:"noticeUserId,omitempty"` //通知用户ID
	NickName     string `json:"nickName,omitempty"`     //用户昵称
	Email        string `json:"email,omitempty"`        //邮箱
	Phone        string `json:"phone,omitempty"`        //手机号
	Avatar       string `json:"avatar,omitempty"`       //头像
}
