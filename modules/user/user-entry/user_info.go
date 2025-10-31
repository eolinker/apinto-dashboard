package user_entry

type UserInfo struct {
	Id           int    `gorm:"column:id;type:INT(10);primary_key;NOT NULL"`
	Sex          int    `gorm:"column:sex;type:INT(11);NOT NULL;comment:性别 0:未知 1:男 2:女"`
	UserName     string `gorm:"column:username;uniqueIndex:user_pk2;type:VARCHAR(36);NOT NULL;comment:用户名"`
	NoticeUserId string `gorm:"column:notice;type:VARCHAR(36);comment:通知用户ID"`
	NickName     string `gorm:"column:nickname;type:VARCHAR(255);comment:昵称"`
	Email        string `gorm:"column:email;type:VARCHAR(255);comment:邮箱"`
	Phone        string `gorm:"column:phone;type:VARCHAR(20);comment:手机号"`
	Avatar       string `gorm:"column:avatar;type:VARCHAR(255);comment:头像"`
	Password     string `gorm:"column:password;not null;type:VARCHAR(32);comment:密码"`
}

func (u *UserInfo) TableName() string {
	return "user"
}

func (u *UserInfo) IdValue() int {
	return u.Id
}
