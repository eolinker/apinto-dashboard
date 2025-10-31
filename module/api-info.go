package apinto_module

import (
	"github.com/eolinker/apinto-dashboard/pm3"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ApiInfoSetHandler(api pm3.Api) gin.HandlerFunc {
	authorityValue := api.Authority.String()
	if authorityValue == "unset" {
		switch api.Method {
		case http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch:
			authorityValue = pm3.Private.String()
		default:
			authorityValue = pm3.Internal.String()
		}
	}

	return func(context *gin.Context) {
		context.Set("apinto-access", api.Access)
		context.Set("apinto-api-authority", authorityValue)
	}
}

type ValueGet interface {
	GetString(string) string
}

func ReadApiInfo(ctx ValueGet) (string, pm3.ApiAuthority) {

	return ctx.GetString("apinto-access"), pm3.ParseApiAuthority(ctx.GetString("apinto-api-authority"))
}
