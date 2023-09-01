/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package module

import (
	"bufio"
	"context"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func (c *Controller) tail(ginCtx *gin.Context) {
	key := ginCtx.Param("key")
	keyData, err := base64.URLEncoding.DecodeString(key)
	if err != nil {
		ginCtx.AbortWithStatus(http.StatusBadRequest)
		return

	}
	url := string(keyData)
	if tokenListContainsValue(ginCtx.Request.Header, "Connection", "Upgrade") {

		var brw *bufio.ReadWriter
		netConn, brw, err := ginCtx.Writer.Hijack()
		if err != nil {
			ginCtx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		go func() {
			defer func() {

				netConn.Close()
			}()
			if brw.Reader.Buffered() > 0 {
				ginCtx.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			upstream, resp, err := DefaultDialer.DialContext(context.Background(), url, ginCtx.Request.Header)
			if err != nil {
				ginCtx.AbortWithStatus(http.StatusInternalServerError)
				return
			}
			defer func() {
				upstream.Close()
			}()
			err = resp.Write(netConn)
			if err != nil {
				return
			}
			go func() {
				io.Copy(upstream, netConn)
			}()
			io.Copy(netConn, upstream)
		}()

	} else {
		ginCtx.AbortWithStatus(http.StatusInternalServerError)
	}
}
