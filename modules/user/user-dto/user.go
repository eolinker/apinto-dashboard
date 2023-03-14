package user_dto

import (
	"github.com/eolinker/apinto-dashboard/access"
)

type SystemModuleItem struct {
	ID       int                  `json:"id"`
	Title    string               `json:"title"`
	Module   string               `json:"module"`
	Access   []*access.AccessItem `json:"access"`
	Children []*SystemModuleItem  `json:"children"`
}

type UserModuleItem struct {
	Id     int      `json:"id"`
	Router string   `json:"router"`
	Title  string   `json:"title"`
	Access []string `json:"access"`
	Parent int      `json:"parent"`
}

type ProxyRoleInfo struct {
	Title  string   `json:"title"`
	Desc   string   `json:"desc"`
	Access []string `json:"access"`
}

type RoleListItem struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	UserNum        int    `json:"user_num"`
	OperateDisable bool   `json:"operate_disable"`
	Type           int    `json:"type"`
}

type RoleOptionItem struct {
	ID             string `json:"id"`
	Title          string `json:"title"`
	OperateDisable bool   `json:"operate_disable"`
}

type BatchUpdateRole struct {
	Ids    []int  `json:"ids"`
	RoleId string `json:"role_id"`
}

type BatchRemoveRole struct {
	Ids    []int  `json:"ids"`
	RoleId string `json:"role_id"`
}

type DelUserReq struct {
	UserIds []int `json:"ids"`
}

type PatchUserReq struct {
	Role   []string `json:"role"`
	Status int      `json:"status"`
}

type SaveUserReq struct {
	Sex          int      `json:"sex"`
	Avatar       string   `json:"avatar"` //头像
	Email        string   `json:"email"`
	Phone        string   `json:"phone"`
	NickName     string   `json:"nick_name"`
	UserName     string   `json:"user_name"`
	Desc         string   `json:"describe"`
	NoticeUserId string   `json:"notice_user_id"`
	RoleIds      []string `json:"role_ids"`
}

type UserInfo struct {
	Id             int      `json:"id"`
	Sex            int      `json:"sex"`
	Avatar         string   `json:"avatar"`
	Email          string   `json:"email"`
	Phone          string   `json:"phone"`
	Status         int      `json:"status"`
	UserName       string   `json:"user_name"`
	NickName       string   `json:"nick_name"`
	LastLogin      string   `json:"last_login"`
	CreateTime     string   `json:"create_time"`
	UpdateTime     string   `json:"update_time"`
	OperateDisable bool     `json:"operate_disable"`
	NoticeUserId   string   `json:"notice_user_id"`
	Desc           string   `json:"desc"`
	Operator       string   `json:"operator"`
	RoleIds        []string `json:"role_ids"`
}

type ResetPasswordReq struct {
	Id       int    `json:"id"`
	Password string `json:"password"`
}

type UpdateMyProfileReq struct {
	NickName     string `json:"nick_name"`
	Email        string `json:"email"`
	Desc         string `json:"desc"`
	NoticeUserId string `json:"notice_user_id"`
}

type UpdateMyPasswordReq struct {
	Old      string `json:"old"`
	Password string `json:"password"`
}
