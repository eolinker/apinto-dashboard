package client

type dataProxy struct {
	Data []byte
}

func (c *dataProxy) UnmarshalJSON(bytes []byte) error {
	c.Data = bytes
	return nil
}

type UserCenterResult struct {
	Success   bool        `json:"success"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	RequestId interface{} `json:"requestId"`
	Data      dataProxy   `json:"data"`
}

type UserInfoReq struct {
	AccountName  string   `json:"accountName,omitempty"`
	Email        string   `json:"email,omitempty"`
	UserName     string   `json:"user_name,omitempty"`
	Phone        string   `json:"phone,omitempty"`
	Id           string   `json:"id,omitempty"`
	UserNickName string   `json:"user_nick_name,omitempty"`
	Keyword      string   `json:"keyword,omitempty"`
	Status       int      `json:"status,omitempty"`
	ExcludeIds   []string `json:"excludeIds,omitempty"`
	Ids          []string `json:"ids,omitempty"`
	GroupIds     string   `json:"groupIds,omitempty"`
	Sort         string   `json:"sort,omitempty"`
	OrderBy      string   `json:"orderBy,omitempty"`
}

type UserInfoRes struct {
	Id            int    `json:"id,omitempty"`
	Sex           int    `json:"sex,omitempty"`
	Password      string `json:"password,omitempty"`
	Avatar        string `json:"avatar,omitempty"`
	Email         string `json:"email,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Status        int    `json:"status,omitempty"`
	GroupId       int    `json:"groupId,omitempty"`
	UserName      string `json:"userName,omitempty"`
	UserNickName  string `json:"userNickName,omitempty"`
	CreateTime    int    `json:"createTime,omitempty"`
	UpdateTime    int    `json:"updateTime,omitempty"`
	UserName1     string `json:"user_name,omitempty"`
	UserNickName1 string `json:"user_nick_name,omitempty"`
	CreateTime1   int    `json:"create_time,omitempty"`
	UpdateTime1   int    `json:"update_time,omitempty"`
}

type SensitiveReq struct {
	Client int    `json:"client,omitempty"`
	UserId string `json:"userId,omitempty"`
	Jwt    string `json:"jwt,omitempty"`
}

type AccessInitReq struct {
	AppId           string   `json:"appId,omitempty"`
	AppName         string   `json:"appName,omitempty"`
	AppKey          string   `json:"appKey,omitempty"`
	AccessTableJson string   `json:"accessTableJson,omitempty"`
	AccessKeys      []string `json:"accessKeys,omitempty"`
	RoleList        []Role   `json:"roleList,omitempty"`
}

type Role struct {
	Id          string   `json:"id,omitempty"`
	Title       string   `json:"title,omitempty"`
	Key         string   `json:"key,omitempty"`
	Description string   `json:"description,omitempty"`
	AppKey      string   `json:"appKey,omitempty"`
	Module      string   `json:"module,omitempty"`
	IsDefault   string   `json:"isDefault,omitempty"`
	AccessList  []string `json:"accessList,omitempty"`
}

type UserAccess struct {
	AccessKeys []string
}

type RoleListReq struct {
	Id     int    `json:"id"`
	Type   int    `json:"type"` //类型（0=内置 1=自定义）
	AppKey string `json:"appKey"`
	Module string `json:"module"` // 默认/
}

type RoleList struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Key         string `json:"key"`
	Description string `json:"description"`
	Type        int    `json:"type"` //类型（0=内置 1=自定义）
	AppId       int    `json:"appId"`
	Module      string `json:"module"`
	IsDefault   int    `json:"isDefault"`
	CreateTime  int    `json:"createTime"`
	UpdateTime  int    `json:"updateTime"`
	AccessKeys  []int  `json:"accessKeys"`
	UserNum     int    `json:"userNum"`
}

type RoleListRes struct {
	RoleList []RoleList `json:"role_list"`
}

type DelUserReq struct {
	UserIds []int `json:"ids"`
}

type UpdateUserPwdReq struct {
	UserId      int    `json:"userId"`
	AccountName string `json:"accountName"`
	IsHash      int    `json:"isHash"`
	OldPassword string `json:"oldPassword"`
	Password    string `json:"password"`
}

type RefreshTokenReq struct {
	RJwt string `json:"rjwt"`
}

type RefreshTokenRes struct {
	Jwt  string `json:"jwt"`
	RJwt string `json:"rjwt"`
}

type ResetUserPwd struct {
	UserId   int    `json:"userId"`
	Password string `json:"password"`
}

type UpdateUserReq struct {
	AccountName  string `json:"accountName,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
	UserNickName string `json:"userNickName,omitempty"`
	Sex          int    `json:"sex,omitempty"`
	Email        string `json:"email,omitempty"`
	UserName     string `json:"userName,omitempty"`
	Phone        string `json:"phone,omitempty"`
	Remark       string `json:"remark,omitempty"`
}

type CreateUserReq struct {
	OperatorId   int      `json:"operatorId"`
	UserNickName string   `json:"userNickName"`
	UserName     string   `json:"userName"`
	Sex          int      `json:"sex"`
	Password     string   `json:"password"`
	Avatar       string   `json:"avatar"`
	Email        string   `json:"email"`
	Phone        string   `json:"phone"`
	Status       int      `json:"status"`
	LastLogin    string   `json:"lastLogin"`
	Remark       string   `json:"remark"`
	AmsInfo      *AmsInfo `json:"amsInfo"`
}
type AmsInfo struct {
	Origin int `json:"origin"` //来源，0表示手动创建，1表示从LDAP导入，2表示oauth登陆
}

type UserAttrResponse struct {
	UserId         string `json:"userId"`
	RegisterWay    string `json:"registerWay"`
	Region         string `json:"region"`
	Apikey         string `json:"apikey"`
	IsActivedEmail string `json:"isActivedEmail"`
	InviteCode     string `json:"inviteCode"`
	InviterId      string `json:"inviterId"`
	AmsVersion     string `json:"amsVersion"`
	CreateTime     string `json:"createTime"`
	UpdateTime     string `json:"updateTime"`
	Ref            string `json:"ref"`
	RefUrl         string `json:"refUrl"`
}
type CreateUserRes struct {
	UserAttrResponse UserAttrResponse `json:"userAttrResponse"`
	Integral         string           `json:"integral"`
	Id               int              `json:"id"`
	Sex              string           `json:"sex"`
	Password         string           `json:"password"`
	Avatar           string           `json:"avatar"`
	Phone            string           `json:"phone"`
	Status           int              `json:"status"`
	UserName         string           `json:"userName"`
	OperatorId       int              `json:"operatorId"`
	UserNickName     string           `json:"userNickName"`
	CreateTime       int              `json:"createTime"`
	UpdateTime       int              `json:"updateTime"`
	UserNameUL       string           `json:"userNameUL"`
	OperatorIdUL     int              `json:"operatorIdUL"`
	UserNickNameUL   string           `json:"userNickNameUL"`
	LastLoginUL      string           `json:"lastLoginUL"`
	CreateTimeUL     int              `json:"createTimeUL"`
	UpdateTimeUL     int              `json:"updateTimeUL"`
}

type LogoutReq struct {
}

type LoginCheck struct {
	Authorization string `json:"authorization"`
}

type UserLoginReq struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Client        int    `json:"client"`
	Type          int    `json:"type"`
	CheckCodeType int    `json:"checkCodeType,omitempty"`
	CheckCode     string `json:"checkCode,omitempty"`
	AppType       int    `json:"appType"`
	Ref           string `json:"ref,omitempty"`
	RefUrl        string `json:"refUrl,omitempty"`
}

type UserLoginRes struct {
	Jwt  string `json:"jwt,omitempty"`
	RJWT string `json:"RJWT,omitempty"`
	Type int    `json:"type,omitempty"`
}
