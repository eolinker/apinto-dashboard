package plugin_client

import (
	"github.com/eolinker/apinto-dashboard/plugin/go-plugin/proto"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
)

type proxy struct {
	client proto.ServiceClient
	id     string
}

func newProxy(client proto.ServiceClient, id string) *proxy {
	return &proxy{client: client, id: id}
}

func (a *proxy) Handle(ginCtx *gin.Context) {
	req := readRequest(ginCtx)
	req.Id = a.id
	resp, err := a.client.Request(ginCtx, req)
	if err != nil {
		log.Errorf(ginCtx.AbortWithError(500, err).Error())
		return
	}
	header := ginCtx.Writer.Header()
	for _, h := range resp.Headers {
		for _, v := range h.Value {
			header.Add(h.Key, v)
		}
	}

	contentType := header.Get("content-type")
	ginCtx.Data(int(resp.Status), contentType, resp.Body)
}
