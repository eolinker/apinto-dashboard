package controller

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/eolinker/apinto-dashboard/access"
	"github.com/eolinker/apinto-dashboard/cache"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/model"
	"github.com/eolinker/apinto-dashboard/service"
	"github.com/eolinker/apinto-dashboard/user_center/client"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"io"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	namespace        = "namespace"
	NamespaceId      = "namespaceId"
	DefaultNamespace = "default"
	UserId           = "userId"
	Authorization    = "Authorization"
	Session          = "Session"
	RAuthorization   = "RAuthorization"
)

var loginError = errors.New("登录已失效，请重新登录")

type IFilterRouter interface {
	MustNamespace(*gin.Context)
	VerifyToken(*gin.Context)
}

var (
	namespaceService service.INamespaceService
	authService      service.IBussinessAuthService
	sessionCache     cache.ISessionCache
	userCenterClient client.IUserCenterClient

	userInfoCache   cache.IUserInfoCache
	roleAccessCache cache.IRoleAccessCache
	userInfoService service.IUserInfoService
)

func init() {

	bean.Autowired(&userInfoService)
	bean.Autowired(&namespaceService)
	bean.Autowired(&sessionCache)
	bean.Autowired(&userCenterClient)
	bean.Autowired(&userInfoCache)
	bean.Autowired(&roleAccessCache)
	bean.Autowired(&authService)
}

func GetNamespaceId(ginCtx *gin.Context) int {
	i := ginCtx.GetInt(NamespaceId)
	if i == 0 {
		i = 1
	}
	return i
}

func LoadNamespaceId(name string) int {
	value, err := namespaceService.GetByName(name)
	if err != nil {
		log.Errorf("filterRouter-loadNamespaceId err=%s name=%s", err.Error(), name)
		return 1
	}
	return value.Id
}

func MustNamespace(context *gin.Context) {

	if namespaceValue := context.Query(namespace); namespaceValue != "" {
		context.Set(namespace, namespaceValue)
		context.Set(NamespaceId, LoadNamespaceId(namespaceValue))
		return
	}
	if namespaceValue := context.GetHeader(namespace); namespaceValue != "" {
		context.Set(namespace, namespaceValue)
		context.Set(NamespaceId, LoadNamespaceId(namespaceValue))
		return
	}
	if namespaceValue, _ := context.Cookie(namespace); namespaceValue != "" {
		context.Set(namespace, namespaceValue)
		context.Set(NamespaceId, LoadNamespaceId(namespaceValue))
		return
	}
	context.Set(namespace, DefaultNamespace)
	context.Set(NamespaceId, LoadNamespaceId(DefaultNamespace))
	//context.JSON(http.StatusOK, dto.NewErrorResult("namespace NotFound"))
	//context.Abort()

}

func VerifyAuth(ginCtx *gin.Context) {
	ctx := context.Background()

	valid, err := authService.CheckCertValid(ctx)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewCertExceedError("授权错误:"+err.Error()))
		ginCtx.Abort()
		return
	}
	if !valid {
		ginCtx.JSON(http.StatusOK, dto.NewCertExceedError("授权已过期"))
		ginCtx.Abort()
		return
	}

}
func VerifyToken(ginCtx *gin.Context) {

	ctx := context.Background()

	session, _ := ginCtx.Cookie(Session)
	if session == "" {
		ginCtx.JSON(http.StatusOK, dto.NewLoginInvalidError(dto.CodeLoginInvalid, loginError.Error()))
		ginCtx.Abort()
		return
	}

	tokens, err := sessionCache.Get(ctx, sessionCache.Key(session))
	if err == redis.Nil || tokens == nil {
		ginCtx.JSON(http.StatusOK, dto.NewLoginInvalidError(dto.CodeLoginInvalid, loginError.Error()))
		ginCtx.Abort()
		return
	}

	token := tokens.Jwt
	rToken := tokens.RJwt

	//1.从ginCtx的header中拿到token，没拿到报错提醒用户重新登录
	verifyToken, err := common.VerifyToken(token)
	if err != nil {
		ginCtx.JSON(http.StatusOK, dto.NewLoginInvalidError(dto.CodeLoginInvalid, loginError.Error()))
		ginCtx.Abort()
		return
	}
	//1.1拿到用户ID和过期时间 过期了重新登录
	claims := verifyToken.Claims.(jwt.MapClaims)
	if err = claims.Valid(); err != nil {
		//续期token
		refreshReq := &client.RefreshTokenReq{RJwt: rToken}
		refreshToken, err := userCenterClient.RefreshToken(refreshReq)
		if err != nil {
			ginCtx.JSON(http.StatusOK, dto.NewLoginInvalidError(dto.CodeLoginInvalid, loginError.Error()))
			ginCtx.Abort()
			return
		}
		ginCtx.Writer.Header().Set(Authorization, refreshToken.Jwt)
		ginCtx.Writer.Header().Set(RAuthorization, refreshToken.RJwt)

		tokens.RJwt = refreshToken.RJwt
		tokens.Jwt = refreshToken.Jwt
		_ = sessionCache.Set(ctx, sessionCache.Key(session), tokens, time.Hour*24*7)

		return
	}

	ginCtx.Writer.Header().Set(Authorization, token)
	ginCtx.Writer.Header().Set(RAuthorization, rToken)

	userId, _ := strconv.Atoi(claims[UserId].(string))

	ginCtx.Set(UserId, userId)
}

func GetUserId(ginCtx *gin.Context) int {
	return ginCtx.GetInt(UserId)
}

func GenAccessHandler(acs ...access.Access) gin.HandlerFunc {

	return func(ginCtx *gin.Context) {
		userId := GetUserId(ginCtx)

		ctx := context.Background()

		var err error
		userInfo, err := userInfoCache.Get(ctx, userInfoCache.Key(userId))
		if err != nil {
			if err == redis.Nil {
				//若缓存没有，则查表，并存入缓存
				user, err := userInfoService.GetUserInfo(context.Background(), userId)
				if err != nil {
					ginCtx.JSON(http.StatusOK, dto.NewNoAccessError(err.Error()))
					ginCtx.Abort()
				}
				userInfo = user.UserInfo
			} else {
				ginCtx.JSON(http.StatusOK, dto.NewNoAccessError(err.Error()))
				ginCtx.Abort()
			}
		}
		roleIds := strings.Split(userInfo.RoleIds, ",")

		for _, role := range roleIds {
			roleAccess, err := roleAccessCache.Get(ctx, roleAccessCache.Key(role))
			if err != nil {
				if err == redis.Nil {
					//若缓存没有，则查表，并存入缓存
					accessIds, err := userInfoService.GetRoleAccessIds(context.Background(), role)
					if err != nil {
						ginCtx.JSON(http.StatusOK, dto.NewNoAccessError(err.Error()))
						ginCtx.Abort()
					}
					accessList := make([]access.Access, 0, len(accessIds))
					for _, accessId := range accessIds {
						id, _ := strconv.Atoi(accessId)
						accessList = append(accessList, access.Access(id))
					}
					roleAccess = model.CreateRoleAccess(accessList...)
					err = roleAccessCache.Set(ctx, roleAccessCache.Key(role), roleAccess, time.Hour)
					if err != nil {
						ginCtx.JSON(http.StatusOK, dto.NewNoAccessError(err.Error()))
						ginCtx.Abort()
					}
				} else {
					ginCtx.JSON(http.StatusOK, dto.NewNoAccessError(err.Error()))
					ginCtx.Abort()
				}
			}
			if roleAccess.IsAccess(acs...) {
				return
			}
		}
		ginCtx.JSON(http.StatusOK, dto.NewNoAccessError("权限不足"))
		ginCtx.Abort()
	}
}

func panicTrace(err interface{}) string {
	buf := new(bytes.Buffer)
	_, _ = fmt.Fprintf(buf, "%v\n", err)
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		_, _ = fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
	}
	return buf.String()
}

func Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("panic %v", panicTrace(err))
			c.JSON(http.StatusInternalServerError, dto.NewErrorResult("服务器内部错误"))
		}
	}()

	c.Next()
}

type ResponseWriterWrapper struct {
	gin.ResponseWriter
	Body *bytes.Buffer // 缓存
}

func (w ResponseWriterWrapper) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w ResponseWriterWrapper) WriteString(s string) (int, error) {
	w.Body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func Logger(ginCtx *gin.Context) {
	if !strings.HasPrefix(ginCtx.Request.URL.Path, "/api/") {
		ginCtx.Next()
		return
	}

	ginCtx.Request = ginCtx.Request.WithContext(context.Background())

	data, err := ginCtx.GetRawData()
	if err == nil {
		ginCtx.Request.Body = io.NopCloser(bytes.NewBuffer(data)) // 关键点
	}

	blw := &ResponseWriterWrapper{Body: bytes.NewBufferString(""), ResponseWriter: ginCtx.Writer}
	ginCtx.Writer = blw

	// 开始时间
	startTime := time.Now()

	// 处理请求
	ginCtx.Next()

	// 结束时间
	endTime := time.Now()

	// 执行时间
	latencyTime := endTime.Sub(startTime)

	// 请求方式
	reqMethod := ginCtx.Request.Method

	// 请求路由
	reqUri := ginCtx.Request.RequestURI

	// 状态码
	statusCode := ginCtx.Writer.Status()

	// 请求IP
	clientIP := ginCtx.ClientIP()
	resBody := blw.Body.String()
	//日志格式
	log.DebugF("| %3d | %13v | %15s | %s | %s | reqBody=%s | resBody=%s ",
		statusCode,
		latencyTime,
		clientIP,
		reqMethod,
		reqUri,
		string(data),
		resBody,
	)
}

type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w CustomResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w CustomResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
