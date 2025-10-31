package middleware

import (
	"github.com/eolinker/apinto-dashboard/pm3"
	"strings"
)

type apiPrefixCheck string

func (a apiPrefixCheck) Check(api pm3.ApiInfo) bool {
	return strings.HasPrefix(api.Path, string(a))
}

var (
	IsApi     = apiPrefixCheck("/api/")
	IsOpenApi = apiPrefixCheck("/api2/")
)
