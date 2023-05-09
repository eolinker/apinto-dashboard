package enum

import "github.com/eolinker/apinto-dashboard/modules/strategy/config"

// Keyword 关键字
var Keyword = map[string]struct{}{config.FilterValuesALL: {}, config.FilterApplication: {}, config.FilterApi: {}, config.FilterPath: {}, config.FilterService: {}, config.FilterMethod: {}, config.FilterIP: {}}
