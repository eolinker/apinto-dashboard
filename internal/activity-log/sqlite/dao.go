package sqlite

import (
	apinto "github.com/eolinker/apinto-dashboard"
)

type ISqliteHandler interface {
	apinto.ActivityLogGetHandler
	apinto.ActivityLogAddHandler
}
