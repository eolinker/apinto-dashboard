package user_entry

import "time"

type UserInfo struct {
	Id            int        `gorm:"column:id"`         //用户ID
	Sex           int        `gorm:"column:sex"`        //0未知 1男2女
	UserName      string     `gorm:"column:username"`   //账号
	Password      string     `gorm:"column:password"`   //密码
	NoticeUserId  string     `gorm:"column:notice""`    //通知用户ID
	NickName      string     `gorm:"column:nickname""`  //用户昵称
	Email         string     `gorm:"column:email"`      //邮箱
	Phone         string     `gorm:"column:phone"`      //手机号
	Avatar        string     `gorm:"column:avatar"`     //头像
	LastLoginTime *time.Time `gorm:"column:login_time"` //最近登录时间
}

func (u *UserInfo) TableName() string {
	return "user"
}

func (u *UserInfo) IdValue() int {
	return u.Id
}
