/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package module

import (
	"encoding/base64"
	"errors"
	"github.com/eolinker/eosc/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (c *Controller) download(ginCtx *gin.Context) {
	key := ginCtx.Param("key")
	keyData, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		ginCtx.AbortWithError(http.StatusBadRequest, errors.New("无效请求"))
		return

	}
	url := string(keyData)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ginCtx.AbortWithError(http.StatusServiceUnavailable, errors.New("无效请求"))
		log.Info("parse url: ", url, " ", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ginCtx.AbortWithError(http.StatusServiceUnavailable, err)
		log.Info("read file form ", url, " ", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		ginCtx.AbortWithError(http.StatusServiceUnavailable, errors.New(resp.Status))
		return
	}

	for k, vs := range resp.Header {
		for _, v := range vs {
			ginCtx.Writer.Header().Add(k, v)
		}
	}

	buf := make([]byte, 4096)
	for {
		read, err := resp.Body.Read(buf)
		if read > 0 {
			ginCtx.Writer.Write(buf[:read])
			ginCtx.Writer.Flush()
		}
		if err != nil {
			break
		}
	}
	ginCtx.Writer.Flush()

}
