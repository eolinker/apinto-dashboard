package users

import (
	"github.com/eolinker/apinto-dashboard/controller"
	"github.com/eolinker/apinto-dashboard/modules/user"
	"github.com/eolinker/eosc/common/bean"
	"github.com/gin-gonic/gin"
)

const (
	userIdKey = "userId"
)

var (
	userInfoService user.IUserInfoService
)

func init() {
	bean.Autowired(&userInfoService)
}
func GetUserId(ginCtx *gin.Context) int {

	id := ginCtx.GetInt(userIdKey)
	if id != 0 {
		return id
	}

	userName := ginCtx.GetString(controller.UserName)
	if userName == "" {
		return 0
	}
	userInfo, err := userInfoService.GetUserInfoByName(ginCtx, userName)
	if err != nil {
		return 0
	}
	ginCtx.Set(userIdKey, userInfo.Id)
	return userInfo.Id

}
