package driver

import (
	"errors"
	"fmt"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/enum"
	api_dto "github.com/eolinker/apinto-dashboard/modules/api/api-dto"
	api_entry "github.com/eolinker/apinto-dashboard/modules/api/api-entry"
	"strings"
)

const apintoRestfulRegexp = "([0-9a-zA-Z-_]+)"

// formatProxyPath 格式化转发路径上，用于转发重写插件正则替换 比如 请求路径 /path/{A}/{B} 原转发路径：/path/{B}  格式化后 新转发路径： /path/$2
func formatProxyPath(requestPath, proxyPath string) string {
	restfulSet := make(map[string]string)
	newProxyPath := proxyPath
	rList := strings.Split(requestPath, "/")
	i := 1
	for _, param := range rList {
		if common.IsRestfulParam(param) {
			restfulSet[param] = fmt.Sprintf("$%d", i)
			i += 1
		}
	}

	for param, order := range restfulSet {
		newProxyPath = strings.ReplaceAll(newProxyPath, param, order)
	}
	return newProxyPath
}

func checkInput(input *api_dto.APIInfo) error {
	input.UUID = strings.TrimSpace(input.UUID)
	if input.UUID != "" {
		err := common.IsMatchString(common.UUIDExp, input.UUID)
		if err != nil {
			return err
		}
	}

	if input.GroupUUID == "" {
		return errors.New("group_uuid can't be nil. ")
	}

	if input.Method == nil {
		input.Method = []string{}
	}

	if input.Match == nil {
		input.Match = []*api_entry.MatchConf{}
	}

	if input.Header == nil {
		input.Header = []*api_entry.ProxyHeader{}
	}

	input.ServiceName = strings.TrimSpace(input.ServiceName)
	if input.ServiceName == "" {
		return errors.New("service_name can't be nil. ")
	}

	var err error
	//格式化,并且校验请求路径
	input.RequestPath = "/" + strings.TrimPrefix(strings.TrimSpace(input.RequestPath), "/")
	if input.RequestPath == "/" {
		return errors.New("request_path can't be / . ")
	}

	if input.RequestPath, err = common.CheckAndFormatPath(input.RequestPath); err != nil {
		return errors.New("request_path is illegal. ")
	}

	requestParamSet := make(map[string]struct{})
	for _, param := range strings.Split(input.RequestPath, "/") {
		if common.IsRestfulParam(param) {
			requestParamSet[param] = struct{}{}
		}
	}

	//格式化转发路径, 并校验转发路径， 若转发路径出现restful参数，则必须对得上请求路径出现过得restful参数
	input.ProxyPath = "/" + strings.TrimPrefix(strings.TrimSpace(input.ProxyPath), "/")
	if input.ProxyPath, err = common.CheckAndFormatPath(input.ProxyPath); err != nil {
		return errors.New("proxy_path is illegal. ")
	}
	for _, param := range strings.Split(input.ProxyPath, "/") {
		if common.IsRestfulParam(param) {
			_, exist := requestParamSet[param]
			if !exist {
				return fmt.Errorf("proxyPath is illegal. restful param %s is invalid. ", param)
			}
		}
	}

	if input.Timeout <= 0 {
		return errors.New("timeout must be greater than 0. ")
	}

	if input.Retry < 0 {
		return errors.New("retry can't be less than 0. ")
	}

	for i, m := range input.Match {
		input.Match[i].Key = strings.TrimSpace(m.Key)
		if input.Match[i].Key == "" {
			return errors.New("Match.Key can't be nil. ")
		}
		switch m.Position {
		case enum.MatchPositionHeader, enum.MatchPositionQuery, enum.MatchPositionCookie:
		default:
			return fmt.Errorf("position %s is illegal. ", m.Position)
		}

		switch m.MatchType {
		case enum.MatchTypeEqual, enum.MatchTypePrefix, enum.MatchTypeSuffix, enum.MatchTypeSubstr, enum.MatchTypeUneuqal, enum.MatchTypeRegexp, enum.MatchTypeRegexpG:
			input.Match[i].Pattern = strings.TrimSpace(m.Pattern)
			if input.Match[i].Pattern == "" {
				return fmt.Errorf("Match.Pattern can't be nil when MatchType is %s. ", m.MatchType)
			}
		case enum.MatchTypeNull, enum.MatchTypeExist, enum.MatchTypeUnexist, enum.MatchTypeAny:
		default:
			return fmt.Errorf("match_type %s is illegal. ", m.MatchType)
		}
	}

	for i, h := range input.Header {
		input.Header[i].Key = strings.TrimSpace(h.Key)
		if input.Header[i].Key == "" {
			return errors.New("Header.Key can't be nil. ")
		}
		switch h.OptType {
		case enum.HeaderOptTypeAdd:
			input.Header[i].Value = strings.TrimSpace(h.Value)
			if input.Header[i].Value == "" {
				return fmt.Errorf("Header.Value can't be nil when OptType is %s. ", h.OptType)
			}
		case enum.HeaderOptTypeDelete:
		default:
			return fmt.Errorf("opt_type %s is illegal. ", h.OptType)
		}
	}

	//设置前端显示所需要的请求路径， 用户输入什么就显示什么
	input.RequestPathLabel = input.RequestPath
	//校验请求路径是不是restful， 如果是则将request_path转换成统一的restful请求路径，如/path/{A} -> /path/{rest}
	if common.IsRestfulPath(input.RequestPath) {
		input.RequestPath = common.ReplaceRestfulPath(input.RequestPath, enum.RestfulLabel)
	}

	if len(input.Hosts) > 0 {
		for _, host := range input.Hosts {
			if !common.IsMatchDomainPort(host) && !common.IsMatchIpPort(host) {
				return fmt.Errorf("host %s is illegal. ", host)
			}
		}
	}

	return nil
}
