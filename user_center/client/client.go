package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common/restful"
	"io"
	"net/http"
	"strings"
)

type IUserCenterClient interface {
	UserInfo(*UserInfoReq) (*UserInfoRes, error)
	DelUser(userIds ...int) error
	CreateUser(req *CreateUserReq) (*int, error)
	UpdateUser(req *UpdateUserReq) error
	UpdateUserPwd(req *UpdateUserPwdReq) error
	ResetUserPwd(req *ResetUserPwd) error
	RefreshToken(req *RefreshTokenReq) (*RefreshTokenRes, error)
	Login(req *UserLoginReq) (*UserLoginRes, int, error)
	LoginCheck(req *LoginCheck) (*bool, error)
	Logout(token string) error
}

type userCenterClient struct {
	userCenterUrl string
	login         restful.Builder[restful.RequestRPC[UserLoginReq, UserLoginRes]]
	refreshToken  restful.Builder[restful.RequestRPC[RefreshTokenReq, RefreshTokenRes]]
	createUser    restful.Builder[restful.RequestRPC[CreateUserReq, int]]
	loginCheck    restful.Builder[restful.RequestRPC[LoginCheck, bool]]
	//userInfo      restful.Builder[restful.RequestRPC[UserInfoReq, []*UserInfoRes]]  数组不适用
	resetUserPwd  restful.Builder[restful.RequestOneWay[ResetUserPwd]]
	updateUserPwd restful.Builder[restful.RequestOneWay[UpdateUserPwdReq]]
	updateUser    restful.Builder[restful.RequestOneWay[UpdateUserReq]]
	delUser       restful.Builder[restful.RequestOneWay[DelUserReq]]
	logout        restful.Builder[restful.RequestOneWay[LogoutReq]]
}

func newIUserCenterClient(url string) IUserCenterClient {
	if strings.TrimSpace(url) == "" {
		panic("user_center_url is null")
	}

	url = strings.TrimSuffix(url, "/")
	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		url = fmt.Sprint("http://", url)
	}

	client := &userCenterClient{userCenterUrl: url}

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	config := restful.BuildConfig(header, nil, url)

	client.login = restful.Rpc[UserLoginReq, UserLoginRes](config, http.MethodPost, "/common/sso/login")
	//client.userInfo = restful.Rpc[UserInfoReq, []*UserInfoRes](config, http.MethodPost, "/user/query")
	client.refreshToken = restful.Rpc[RefreshTokenReq, RefreshTokenRes](config, http.MethodPost, "/common/sso/refresh")
	client.createUser = restful.Rpc[CreateUserReq, int](config, http.MethodPost, "/user/save")
	client.loginCheck = restful.Rpc[LoginCheck, bool](config, http.MethodPost, "/common/sso/login/check")
	client.resetUserPwd = restful.OneWay[ResetUserPwd](config, http.MethodPost, "/user/confirmUser")
	client.updateUserPwd = restful.OneWay[UpdateUserPwdReq](config, http.MethodPost, "/user/changePassword")
	client.updateUser = restful.OneWay[UpdateUserReq](config, http.MethodPost, "/user/changeUser")
	client.delUser = restful.OneWay[DelUserReq](config, http.MethodPost, "/user/soft-delete")
	client.logout = restful.OneWay[LogoutReq](config, http.MethodPost, "/common/sso/logout")

	return client
}

// newRequest body = strings.NewReader(string(by))
func newRequest(method, url string, body io.Reader) (by []byte, err error) {
	method = strings.ToUpper(method)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, body)
	req.Header.Add("Content-Type", "application/json;charset=UTF-8")
	defer req.Body.Close()
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (u *userCenterClient) Login(req *UserLoginReq) (*UserLoginRes, int, error) {
	res, err := u.login.Build().Request(req)
	if err != nil {
		return nil, -1, err
	}

	if res.Success {
		return res.Data, 0, nil
	}

	return nil, res.Code, errors.New(res.Message)
}

func (u *userCenterClient) Logout(token string) error {
	req := &LogoutReq{}
	res, err := u.logout.Build().Header("Authorization", token).Request(req)
	if err != nil {
		return err
	}

	if res.Success {
		return nil
	}

	return errors.New(res.Message)
}

func (u *userCenterClient) RefreshToken(req *RefreshTokenReq) (*RefreshTokenRes, error) {
	res, err := u.refreshToken.Build().Request(req)
	if err != nil {
		return nil, err
	}

	if res.Success {
		return res.Data, nil
	}

	return nil, errors.New(res.Message)

	//infoReqBytes, _ := json.Marshal(req)
	//
	//resBytes, err := newRequest(http.MethodPost, u.userCenterUrl+"/common/sso/refresh", strings.NewReader(string(infoReqBytes)))
	//if err != nil {
	//	return nil, err
	//}
	//
	//result := &UserCenterResult{}
	//err = json.Unmarshal(resBytes, result)
	//if err != nil {
	//	return nil, err
	//}
	//
	////解析用户信息
	//refreshTokenRes := new(RefreshTokenRes)
	//err = json.Unmarshal(result.Data.Data, refreshTokenRes)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return refreshTokenRes, nil
}

func (u *userCenterClient) UserInfo(userInfo *UserInfoReq) (*UserInfoRes, error) {

	infoReqBytes, _ := json.Marshal(userInfo)

	resBytes, err := newRequest(http.MethodPost, u.userCenterUrl+"/user/query", strings.NewReader(string(infoReqBytes)))
	if err != nil {
		return nil, err
	}

	result := &UserCenterResult{}
	err = json.Unmarshal(resBytes, result)
	if err != nil {
		return nil, err
	}

	//解析用户信息
	userInfosRes := make([]*UserInfoRes, 0)
	err = json.Unmarshal(result.Data.Data, &userInfosRes)
	if err != nil {
		return nil, err
	}

	if len(userInfosRes) > 0 {
		return userInfosRes[0], nil
	}
	return nil, nil
}

func (u *userCenterClient) DelUser(userIds ...int) error {
	req := &DelUserReq{UserIds: userIds}

	res, err := u.delUser.Build().Request(req)
	if err != nil {
		return err
	}
	if res.Success {
		return nil
	}
	return errors.New(res.Message)

	//reqBytes, _ := json.Marshal(req)
	//
	//resBytes, err := newRequest(http.MethodPost, u.userCenterUrl+"/user/soft-delete", strings.NewReader(string(reqBytes)))
	//if err != nil {
	//	return err
	//}
	//
	//res := &UserCenterResult{}
	//err = json.Unmarshal(resBytes, res)
	//if err != nil {
	//	return err
	//}
	//
	//if res.Success {
	//	return nil
	//}
	//return errors.New(res.Message)
}

func (u *userCenterClient) UpdateUserPwd(req *UpdateUserPwdReq) error {
	res, err := u.updateUserPwd.Build().Request(req)
	if err != nil {
		return err
	}

	if res.Success {
		return nil
	}
	return errors.New(res.Message)

	//reqBytes, _ := json.Marshal(req)

	//resBytes, err := newRequest(http.MethodPost, u.userCenterUrl+"/common/console/user/password/update", strings.NewReader(string(reqBytes)))
	//if err != nil {
	//	return err
	//}
	//
	//res := &UserCenterResult{}
	//err = json.Unmarshal(resBytes, res)
	//if err != nil {
	//	return err
	//}
	//
	//if res.Success {
	//	return nil
	//}
	//return errors.New(res.Message)
}

func (u *userCenterClient) ResetUserPwd(req *ResetUserPwd) error {

	res, err := u.resetUserPwd.Build().Request(req)
	if err != nil {
		return err
	}

	if res.Success {
		return nil
	}
	return errors.New(res.Message)

	//reqBytes, _ := json.Marshal(req)
	//
	//resBytes, err := newRequest(http.MethodPost, u.userCenterUrl+"/user/confirmUser", strings.NewReader(string(reqBytes)))
	//if err != nil {
	//	return err
	//}
	//
	//res := &UserCenterResult{}
	//err = json.Unmarshal(resBytes, res)
	//if err != nil {
	//	return err
	//}
	//
	//if res.Success {
	//	return nil
	//}
	//return errors.New(res.Message)
}

func (u *userCenterClient) LoginCheck(req *LoginCheck) (*bool, error) {
	res, err := u.loginCheck.Build().Header("Authorization", req.Authorization).Request(req)
	if err != nil {
		return nil, err
	}

	if res.Success {
		return res.Data, nil
	}

	return nil, errors.New(res.Message)
}

func (u *userCenterClient) CreateUser(req *CreateUserReq) (*int, error) {
	req.AmsInfo = new(AmsInfo)
	req.AmsInfo.Origin = 0
	res, err := u.createUser.Build().Request(req)
	if err != nil {
		return nil, err
	}

	if res.Success {
		return res.Data, nil
	}

	return nil, errors.New(res.Message)

	//createUserReqBytes, _ := json.Marshal(req)
	//
	//createUserResBytes, err := newRequest(http.MethodPost, u.userCenterUrl+"/common/console/user/save", strings.NewReader(string(createUserReqBytes)))
	//if err != nil {
	//	return nil, err
	//}
	//res := &UserCenterResult{}
	//err = json.Unmarshal(createUserResBytes, res)
	//if err != nil {
	//	return nil, err
	//}
	//
	//if res.Success {
	//	createUserRes := &CreateUserRes{}
	//	err = json.Unmarshal(res.Data.Data, createUserRes)
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	return createUserRes, nil
	//}
	//return nil, errors.New(res.Message)

}

func (u *userCenterClient) UpdateUser(req *UpdateUserReq) error {

	res, err := u.updateUser.Build().Request(req)
	if err != nil {
		return err
	}

	if res.Success {
		return nil
	}

	return errors.New(res.Message)

	//reqBytes, _ := json.Marshal(req)
	//
	//resBytes, err := newRequest(http.MethodPost, u.userCenterUrl+"/user/changeUser", strings.NewReader(string(reqBytes)))
	//if err != nil {
	//	return err
	//}
	//res := &UserCenterResult{}
	//err = json.Unmarshal(resBytes, res)
	//if err != nil {
	//	return err
	//}
	//if res.Success {
	//	return nil
	//}
	//
	//return errors.New(res.Message)
}
