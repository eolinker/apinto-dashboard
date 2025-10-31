package cache

import (
	"fmt"

	"strings"

	"github.com/eolinker/eosc/common/bean"
	"github.com/redis/go-redis/v9"
)

var (
	client    redis.UniversalClient
	namespace string
)

func InitCache(c redis.UniversalClient, prefix string) {
	namespace = prefix
	if namespace == "" {
		namespace = "apinto"
	}
	namespace = fmt.Sprint(strings.Trim(namespace, ":"), ":")
	client = c
	iCommonCache := newCommonCache(client, namespace)
	bean.Injection(&iCommonCache)
	lockerManger := newManager(client, namespace)
	bean.Injection(&lockerManger)

}
