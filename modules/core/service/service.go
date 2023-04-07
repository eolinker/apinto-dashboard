package service

import (
	"github.com/eolinker/apinto-dashboard/modules/core"
	"net/http"
	"sync"
	"sync/atomic"
)

var (
	_ core.ICore = (*coreService)(nil)
)

type coreService struct {
	handlerPointer      atomic.Pointer[http.Handler]
	localVersionPointer atomic.Pointer[string]
	lock                sync.Mutex
}

func (c *coreService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := c.handlerPointer.Load()
	if handler == nil || (*handler) == nil {
		http.NotFound(w, r)
		return
	}
	(*handler).ServeHTTP(w, r)
}

func (c *coreService) ReloadModule(version string) error {

	localVersion := c.localVersionPointer.Swap(&version)

	if localVersion != nil && (*localVersion) == version {
		return nil
	}
	c.lock.Lock()
	defer c.lock.Unlock()
	localVersion = c.localVersionPointer.Load()
	if localVersion != nil && (*localVersion) == version {
		// todo load module
		// todo load middleware

	}
	return nil
}

func NewService() *coreService {
	return &coreService{}
}
