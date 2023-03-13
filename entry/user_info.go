package entry

import "time"

type UserInfo struct {
	Id            int        `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`  //用户ID
	Sex           int        `gorm:"type:int(11);size:11;not null;default:0;column:sex;comment:0未知1男2女"`               //0未知 1男2女
	UserName      string     `gorm:"size:36;not null;index:user_name;column:user_name;comment:账号名"`                    //账号
	NoticeUserId  string     `gorm:"size:36;not null;column:notice_user_id;comment:通知用户ID"`                            //通知用户ID
	NickName      string     `gorm:"size:36;not null;column:nick_name;comment:昵称"`                                     //用户昵称
	Email         string     `gorm:"size:255;default:null;column:email;comment:邮箱地址"`                                  //邮箱
	Phone         string     `gorm:"size:20;default:null;column:phone;comment:手机号"`                                    //手机号
	Avatar        string     `gorm:"size:255;default:null;column:avatar;comment:头像"`                                   //头像
	Remark        string     `gorm:"size:255;default:null;column:remark;comment:备注"`                                   //备注
	RoleIds       string     `gorm:"size:255;default:null;column:role_ids;comment:角色ID数组"`                             //角色
	Status        int        `gorm:"type:int(11);size:11;not null;default:2;column:status;comment:启用状态，1不启用，2启用"`      //启用状态，1不启用，2启用
	IsDelete      bool       `gorm:"type:tinyint(1);size:1;default:0;is_delete;comment:是否删除"`                          //是否删除
	Operator      int        `gorm:"type:int(11);size:11;default:null;column:operator;comment:操作人"`                    //操作人
	FlushTime     time.Time  `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:flush_time;comment:刷新时间"` //上次刷新时间
	CreateTime    time.Time  `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP;column:create_time;comment:创建时间"`
	UpdateTime    time.Time  `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:update_time;comment:修改时间"`
	LastLoginTime *time.Time `gorm:"type:timestamp;NOT NULL;DEFAULT:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP;column:last_login_time;comment:最近登录时间"` //最近登录时间
}

func (u *UserInfo) TableName() string {
	return "user_info"
}

func (u *UserInfo) IdValue() int {
	return u.Id
}
