package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	expires      = time.Hour * 24 * 7
	cacheControl = fmt.Sprintf("public, max-age=%d", 3600*24*7)
)

func addExpires(ginCtx *gin.Context) {
	ginCtx.Header("Expires", time.Now().Add(expires).UTC().Format(http.TimeFormat))
	ginCtx.Header("Cache-Control", cacheControl)
}
