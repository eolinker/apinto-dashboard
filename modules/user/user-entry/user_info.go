package user_entry

import "time"

type UserInfo struct {
	Id            int        `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`                                //用户ID
	Sex           int        `gorm:"type:int(11);size:11;not null;default:0;column:sex;comment:0未知1男2女"`                                             //0未知 1男2女
	UserName      string     `gorm:"size:36;not null;index:username;dbUniqueIndex:username;column:username;comment:账号名"`                             //账号
	NoticeUserId  string     `gorm:"size:36;not null;column:notice;comment:通知用户ID"`                                                                  //通知用户ID
	NickName      string     `gorm:"size:36;not null;column:nickname;comment:昵称"`                                                                    //用户昵称
	Email         string     `gorm:"size:255;default:null;column:email;comment:邮箱地址"`                                                                //邮箱
	Phone         string     `gorm:"size:20;default:null;column:phone;comment:手机号"`                                                                  //手机号
	Avatar        string     `gorm:"size:255;default:null;column:avatar;comment:头像"`                                                                 //头像
	LastLoginTime *time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:login_time;comment:最近登录时间"` //最近登录时间
}

func (u *UserInfo) TableName() string {
	return "user"
}

func (u *UserInfo) IdValue() int {
	return u.Id
}
