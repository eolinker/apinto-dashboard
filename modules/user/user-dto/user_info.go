package user_dto

type UserInfo struct {
	Id           int    `json:"id"`                     //用户ID
	Sex          int    `json:"sex,omitempty"`          //0未知 1男2女
	UserName     string `json:"userName,omitempty"`     //账号
	NoticeUserId string `json:"noticeUserId,omitempty"` //通知用户ID
	NickName     string `json:"nickName,omitempty"`     //用户昵称
	Email        string `json:"email,omitempty"`        //邮箱
	Phone        string `json:"phone,omitempty"`        //手机号
	Avatar       string `json:"avatar,omitempty"`       //头像
	LastLogin    string `json:"last_login,omitempty"`   //最近登录时间

}
type UserEnum struct {
	Sex          int    `json:"sex,omitempty"`          //0未知 1男2女
	UserName     string `json:"userName,omitempty"`     //账号
	NoticeUserId string `json:"noticeUserId,omitempty"` //通知用户ID
	NickName     string `json:"nickName,omitempty"`     //用户昵称
	Email        string `json:"email,omitempty"`        //邮箱
	Phone        string `json:"phone,omitempty"`        //手机号
	Avatar       string `json:"avatar,omitempty"`       //头像
}

type ResetPasswordReq struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
}

type UpdateMyProfileReq struct {
	NickName string `json:"nick_name"`
	Email    string `json:"email"`
	//Desc         string `json:"desc"`
	NoticeUserId string `json:"notice_user_id"`
}

type UpdateMyPasswordReq struct {
	Old      string `json:"old"`
	Password string `json:"password"`
}
