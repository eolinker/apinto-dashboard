package namespace_controller

import (
	namespace2 "github.com/eolinker/apinto-dashboard/modules/namespace"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
)

var (
	namespaceService namespace2.INamespaceService
)

const (
	namespace        = "namespace"
	NamespaceId      = "namespaceId"
	DefaultNamespace = "default"
)

func init() {
	bean.Autowired(&namespaceService)
}
func LoadNamespaceId(name string) int {
	value, err := namespaceService.GetByName(name)
	if err != nil {
		log.Errorf("filterRouter-loadNamespaceId err=%s name=%s", err.Error(), name)
		return 1
	}
	return value.Id
}
func GetNamespaceId(ginCtx *gin.Context) int {
	i := ginCtx.GetInt(NamespaceId)
	if i == 0 {
		i = 1
	}
	return i
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
