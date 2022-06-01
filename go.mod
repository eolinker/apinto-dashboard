module github.com/eolinker/apinto-dashboard

go 1.17

require (
	github.com/eolinker/eosc v0.4.2
	github.com/go-basic/uuid v1.0.0
	github.com/julienschmidt/httprouter v1.3.0
	gopkg.in/yaml.v2 v2.4.0
	github.com/mattn/go-sqlite3 v1.14.12
)
replace (
	github.com/eolinker/eosc  => ../eosc
)
