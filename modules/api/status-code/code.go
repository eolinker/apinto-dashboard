package status_code

const (
	ApiModuleBaseErrCode = 101000
	ApiConfigBindErr     = ApiModuleBaseErrCode + iota
	ApiSchemeNotExist
	ApiConfigCheckErr
	ApiRouterReduplicatedErr
	ApiGroupNotExistErr
	ApiServiceNotExistErr
	ApiTemplateNotExistErr
	ApiDataBaseErr
)
