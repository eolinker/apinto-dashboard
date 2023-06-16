/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package module

import (
	"github.com/gin-gonic/gin"
)

var (
//	dialer = Dialer{
//		NetDial:          nil,
//		HandshakeTimeout: 0,
//	}
)

func (c *Controller) tail(ginCtx *gin.Context) {

	//if !tokenListContainsValue(ginCtx.Request.Header, "Connection", "Upgrade") {
	//
	//	var brw *bufio.ReadWriter
	//	netConn, brw, err := ginCtx.Writer.Hijack()
	//	if err != nil {
	//		ginCtx.AbortWithStatus(http.StatusInternalServerError)
	//		return
	//	}
	//	defer func() {
	//
	//		netConn.Close()
	//	}()
	//	if brw.Reader.Buffered() > 0 {
	//
	//		return
	//	}
	//
	//
	//	upstream, resp, err := dialer.DialContext(request)
	//	if err != nil {
	//		ginCtx.AbortWithStatus(http.StatusInternalServerError)
	//		return
	//	}
	//	defer upstream.Close()
	//	err = resp.Write(netConn)
	//	if err != nil {
	//		return
	//	}
	//	go func() {
	//		io.Copy(netConn, upstream)
	//	}()
	//	io.Copy(upstream, netConn)
	//} else {
	//	ginCtx.AbortWithStatus(http.StatusInternalServerError)
	//}
}
