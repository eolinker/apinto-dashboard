package v1

// node model

type Node struct {
	Id     string   `json:"id"`
	Name   string   `json:"name"`
	Peer   []string `yaml:"peer"`
	Admin  []string `json:"admin"`
	Server []string `json:"server"`
	Leader bool     `json:"leader"`
}

type ClusterInfo struct {
	Cluster string  `yaml:"cluster"`
	Nodes   []*Node `yaml:"nodes"`
}

type WorkerInfo struct {
	Id          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Profession  string `json:"profession,omitempty"`
	Driver      string `json:"driver,omitempty"`
	Description string `json:"description,omitempty"`
	Update      string `json:"update,omitempty"`
	Create      string `json:"create,omitempty"`
	Disable     bool   `json:"disable"`
}

type RouterInfo struct {
	Create      string   `json:"create"`
	Description string   `json:"description"`
	Disable     bool     `json:"disable"`
	Driver      string   `json:"driver"`
	Host        []string `json:"host"`
	ID          string   `json:"id"`
	Listen      int      `json:"listen"`
	Name        string   `json:"name"`
	Profession  string   `json:"profession"`
	Target      string   `json:"target"`
	Update      string   `json:"update"`
}

// auth model
type AuthConfig struct {
	HideCredentials bool        `json:"hide_credentials,omitempty"`
	User            interface{} `json:"user"` //[]AkSkUser,[]ApiKeyUser,BasicUser[]
	Name            string      `json:"name"`
	Driver          string      `json:"driver"`
	Description     string      `json:"description"`
	*JwtConfig
}

type AkSkUser struct {
	Ak     string            `json:"ak"`
	Sk     string            `json:"sk"`
	Labels map[string]string `json:"labels,omitempty"`
	Expire int64             `json:"expire"`
}

type ApiKeyUser struct {
	Apikey string            `json:"apikey"`
	Labels map[string]string `json:"labels,omitempty"`
	Expire int64             `json:"expire"`
}

type BasicUser struct {
	Username string            `json:"username"`
	Password string            `json:"password"`
	Labels   map[string]string `json:"labels,omitempty"`
	Expire   int64             `json:"expire"` //过期时间 毫秒时间戳
}

type JwtConfig struct {
	Credentials       []*JwtCredentials `json:"credentials,omitempty"`
	SignatureIsBase64 bool              `json:"signature_is_base64,omitempty"` //是否 base64加密
	ClaimsToVerify    []string          `json:"claims_to_verify,omitempty"`    //校验字段
}

type JwtCredentials struct {
	Iss          string            `json:"iss"`            //证书签发者
	Secret       string            `json:"secret"`         //密钥
	RsaPublicKey string            `json:"rsa_public_key"` //公钥
	Algorithm    string            `json:"algorithm"`      //签名算法
	Labels       map[string]string `json:"labels,omitempty"`
}

//discovery model

type DiscoveryConfig struct {
	//Scheme      string      `json:"scheme"`
	Name        string      `json:"name"`
	Driver      string      `json:"driver"`
	Description string      `json:"description"`
	Config      interface{} `json:"config"` //StaticConfig ConsulConfig EurekaConfig NacosConfig
	*StaticHealth
}

type StaticHealth struct {
	HealthOn bool    `json:"health_on"`
	Health   *Health `json:"health,omitempty"`
}

type Health struct {
	Scheme      string `json:"scheme"`
	Method      string `json:"method"`
	Url         string `json:"url"`
	SuccessCode int    `json:"success_code"`
	Period      int    `json:"period"`
	Timeout     int    `json:"timeout"`
}

type StaticConfig struct {
	Address []string          `json:"address"`
	Params  map[string]string `json:"params,omitempty"`
}

type ConsulConfig struct {
	Address []string          `json:"address"`
	Params  map[string]string `json:"params,omitempty"`
}

type EurekaConfig struct {
	Address []string          `json:"address"`
	Params  map[string]string `json:"params,omitempty"`
}

type NacosConfig struct {
	Address []string          `json:"address"`
	Params  map[string]string `json:"params,omitempty"`
}

// RateLimitingConfig 限流策略
type RateLimitingConfig struct {
	Name        string              `json:"name"`
	Stop        bool                `json:"stop"`
	Description string              `json:"description"`
	Priority    int                 `json:"priority"`
	Filter      map[string][]string `json:"filter"`
	Limiting    LimitConf           `json:"limiting"`
}

type LimitConf struct {
	Metrics []string `json:"metrics"`
	Query   Limit    `json:"query"`
	Traffic Limit    `json:"traffic"`
}

type Limit struct {
	Second int `json:"second"`
	Minute int `json:"minute"`
	Hour   int `json:"hour"`
}

//output model

type OutputConfig struct {
	Name        string              `json:"name,omitempty"`
	Driver      string              `json:"driver,omitempty"`
	Type        string              `json:"type,omitempty"` //line,json
	Formatter   map[string][]string `json:"formatter,omitempty"`
	Description string              `json:"description,omitempty"` //描述
}

type FileOutPut struct {
	OutputConfig
	Dir    string `json:"dir,omitempty"`    //目录地址
	File   string `json:"file,omitempty"`   //文件名
	Period string `json:"period,omitempty"` //day,
	Expire int    `json:"expire,omitempty"` //文件保存时间 天
}

type HttpOutPut struct {
	OutputConfig
	Headers map[string]string `json:"headers,omitempty"` //请求头
	Url     string            `json:"url,omitempty"`
	Method  string            `json:"method,omitempty"` //GET POST PUT
}

type KafkaOutPut struct {
	OutputConfig
	Topic         string `json:"topic"`
	Address       string `json:"address"` //请求地址
	Timeout       int    `json:"timeout"`
	Version       string `json:"version"`
	PartitionType string `json:"partition_type"`
	Partition     int    `json:"partition"`
	PartitionKey  string `json:"partition_key"`
}

type NsqdOutPut struct {
	OutputConfig
	Topic      string   `json:"topic,omitempty"`
	Address    []string `json:"address,omitempty"`     //请求地址
	AuthSecret string   `json:"auth_secret,omitempty"` //鉴权secret
}

type SyslogOutPut struct {
	OutputConfig
	Network string `json:"network"` //网络协议 tcp,udp,unix
	Address string `json:"address"` //地址
	Level   string `json:"level"`   //日志等级 error,warn,info,debug,trace
}

//plugin model

type IPluginInfo interface {
	PluginId() string
}

type Plugin struct {
	Disable bool        `json:"disable"` //是否禁用
	Config  interface{} `json:"config"`  //IPluginInfo
}

type GlobalPlugin struct {
	Config     interface{} `json:"config,omitempty"` //Plugin***Config
	Id         string      `json:"id"`
	InitConfig interface{} `json:"init_config,omitempty"`
	Name       string      `json:"name"`   //名称
	Status     string      `json:"status"` //enable,disable,global
	Rely       string      `json:"rely"`   //依赖哪个插件
}

type PluginAuthConfig struct {
	Auth []string `json:"auth,omitempty"`
}

func (PluginAuthConfig) PluginId() string {
	return "eolinker.com:apinto:auth"
}

//type PluginRateLimitingConfig struct {
//	Second           int    `json:"second,omitempty"`             //每秒请求限制
//	Minute           int    `json:"minute,omitempty"`             //每分钟请求限制
//	Hour             int    `json:"hour,omitempty"`               //每小时请求限制
//	Day              int    `json:"day,omitempty"`                //每天请求限制
//	HideClientHeader bool   `json:"hide_client_header,omitempty"` //是否隐藏流控信息
//	ResponseType     string `json:"response_type,omitempty"`      //报错格式 json,text
//}
//
//func (PluginRateLimitingConfig) PluginId() string {
//	return "eolinker.com:apinto:rate_limiting"
//}

type PluginProxyRewriteConfig struct {
	Headers  map[string]string `json:"headers,omitempty"` //请求头部
	Host     string            `json:"host,omitempty"`
	Uri      string            `json:"uri,omitempty"`
	RegexUri []string          `json:"regex_uri,omitempty"` //正则替换路径
}

func (PluginProxyRewriteConfig) PluginId() string {
	return "eolinker.com:apinto:proxy_rewrite"
}

type PluginProxyRewriteV2Config struct {
	PathType    string            `json:"path_type"`
	StaticPath  string            `json:"static_path,omitempty"` //path_type=static启用
	PrefixPath  []*PrefixPath     `json:"prefix_path,omitempty"` //path_type=prefix 启用
	RegexPath   []*RegexPath      `json:"regex_path,omitempty"`  //path_type=regex 启用
	NotMatchErr bool              `json:"not_match_err"`
	HostRewrite bool              `json:"host_rewrite"`
	Host        string            `json:"host,omitempty"`
	Headers     map[string]string `json:"headers"`
}

type RegexPath struct {
	RegexPathMatch   string `json:"regex_path_match"`
	RegexPathReplace string `json:"regex_path_replace"`
}

type PrefixPath struct {
	PrefixPathMatch   string `json:"prefix_path_match"`
	PrefixPathReplace string `json:"prefix_path_replace"`
}

func (PluginProxyRewriteV2Config) PluginId() string {
	return "eolinker.com:apinto:proxy_rewrite_v2"
}

type PluginResponseRewriteConfig struct {
	Body       string            `json:"body"`        //响应内容
	BodyBase64 bool              `json:"body_base64"` //是否base64加密
	Headers    map[string]string `json:"headers"`     //响应头
	Match      struct {
		Code []int `json:"code"`
	} `json:"match"`
	StatusCode int `json:"status_code"` //响应状态码
}

func (PluginResponseRewriteConfig) PluginId() string {
	return "eolinker.com:apinto:response_rewrite"
}

type PluginCircuitBreakerConfig struct {
	MatchCodes      string            `json:"match_codes,omitempty"`      //匹配状态码
	MonitorPeriod   int               `json:"monitor_period,omitempty"`   //监控期
	MinimumRequests int               `json:"minimum_requests,omitempty"` //最低熔断阀值，达到熔断状态的最少请求次数
	FailurePercent  int               `json:"failure_percent,omitempty"`  //监控期内的请求错误率
	BreakPeriod     int               `json:"break_period,omitempty"`     //熔断期
	SuccessCounts   int               `json:"success_counts,omitempty"`   //连续请求成功次数，半开放状态下请求成功次数达到后会转变成健康状态
	BreakerCode     int               `json:"breaker_code,omitempty"`     //熔断状态下返回的响应状态码
	Headers         map[string]string `json:"headers,omitempty"`          //熔断状态下新增的返回头部值
	Body            string            `json:"body,omitempty"`             //熔断状态下的返回响应体
}

func (PluginCircuitBreakerConfig) PluginId() string {
	return "eolinker.com:apinto:circuit_breaker"
}

type PluginIpRestrictionConfig struct {
	IpListType  string   `json:"ip_list_type"`            //white,black
	IpWhiteList []string `json:"ip_white_list,omitempty"` //IpListType为white必填
	IpBlackList []string `json:"ip_black_list,omitempty"` //IpListType为black必填
}

func (PluginIpRestrictionConfig) PluginId() string {
	return "eolinker.com:apinto:ip_restriction"
}

type PluginGZipConfig struct {
	Types     []string `json:"types"`      //content-type列表
	MinLength int      `json:"min_length"` //待压缩内容的最小长度
	Vary      bool     `json:"vary"`       //是否加上Vary头部
}

func (PluginGZipConfig) PluginId() string {
	return "eolinker.com:apinto:gzip"
}

type PluginCorsConfig struct {
	AllowOrigins     string `json:"allow_origins,omitempty"`     //允许跨域访问的Origin
	AllowMethods     string `json:"allow_methods,omitempty"`     //允许通过的请求方式
	AllowCredentials bool   `json:"allow_credentials,omitempty"` //请求中是否携带cookie
	AllowHeaders     string `json:"allow_headers,omitempty"`     //允许跨域访问时请求方携带的非CORS规范以外的Header
	ExposeHeaders    string `json:"expose_headers,omitempty"`    //允许跨域访问时响应方携带的非CORS规范以外的Header
	MaxAge           int    `json:"max_age,omitempty"`           //浏览器缓存CORS结果的最大时间
}

func (PluginCorsConfig) PluginId() string {
	return "eolinker.com:apinto:cors"
}

type PluginTransformerConfig struct {
	Params []struct {
		Name          string `json:"name"`           //待映射参数名称
		Position      string `json:"position"`       //待映射参数所在位置
		ProxyName     string `json:"proxy_name"`     //目标参数名称
		ProxyPosition string `json:"proxy_position"` //目标参数所在位置 可选参数 header,query,body
		Required      bool   `json:"required"`       //待映射参数是否必含
	} `json:"params"`
	Remove    bool   `json:"remove"`
	ErrorType string `json:"error_type"` //可选参数 json,text
}

func (PluginTransformerConfig) PluginId() string {
	return "eolinker.com:apinto:params_transformer"
}

type PluginAccessLogConfig struct {
	Output []interface{} `json:"output"`
}

func (PluginAccessLogConfig) PluginId() string {
	return "eolinker.com:apinto:access_log"
}

type PluginExtraParamsConfig struct {
	Params []struct {
		Name     string `json:"name"`     //参数名
		Position string `json:"position"` //参数位置
		Value    string `json:"value"`    //参数值
		Conflict string `json:"conflict"` //可选参数origin,convert,error
	} `json:"params"`
	ErrorType string `json:"error_type"` //可选参数 json,text
}

func (PluginExtraParamsConfig) PluginId() string {
	return "eolinker.com:apinto:extra_params"
}

// router model

type RouterConfig struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description"` //描述
	Driver      string `json:"driver"`      //http
	Append      map[string]interface{}
}

type RouterRule struct {
	Type  string `json:"type"`  //label:"类型" enum:"header,query,cookie"
	Name  string `json:"name"`  //label:"参数名"
	Value string `json:"value"` //label:"值规"
}

// ApplicationConfig 应用
type ApplicationConfig struct {
	Name        string                  `json:"name,omitempty"`
	Driver      string                  `json:"driver,omitempty"`
	Auth        []ApplicationAuth       `json:"auth"`        //用户列表
	Disable     bool                    `json:"disable"`     //禁用
	Description string                  `json:"description"` //描述
	Labels      map[string]string       `json:"labels,omitempty"`
	Additional  []ApplicationAdditional `json:"additional,omitempty"` //额外参数
	Anonymous   bool                    `json:"anonymous"`            //是否匿名
}

type ApplicationAdditional struct {
	Key      string `json:"key,omitempty"`
	Value    string `json:"value,omitempty"`
	Position string `json:"position,omitempty"` //header,query,body
}

type ApplicationAuth struct {
	Config    interface{}           `json:"config"` //ApplicationAuthBasicConfig ApplicationAuthApiKeyConfig ApplicationAuthAkSkConfig ApplicationAuthJwtConfig
	Type      string                `json:"type"`   //basic apikey  ak/sk jwt
	Users     []ApplicationAuthUser `json:"users"`
	Position  string                `json:"position"`
	TokenName string                `json:"token_name"`
}
type ApplicationAuthUser struct {
	Expire         int64             `json:"expire"`           //秒时间戳
	Labels         map[string]string `json:"labels,omitempty"` //用户标签
	Pattern        interface{}       `json:"pattern"`          //用户匹配规则 ApplicationAuthUserBasicPattern ApplicationAuthUserApiKeyPattern ApplicationAuthUserAkSkPattern ApplicationAuthUserJwtPattern
	HideCredential bool              `json:"hide_credential"`  //转发时是否隐藏鉴权信息
}

type ApplicationAuthUserBasicPattern struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ApplicationAuthUserApiKeyPattern struct {
	Apikey string `json:"apikey"`
}

type ApplicationAuthUserAkSkPattern struct {
	Ak string `json:"ak"`
	Sk string `json:"sk"`
}
type ApplicationAuthUserJwtPattern struct {
	UserName string `json:"username"`
}

type ApplicationAuthBasicConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Expire   int64  `json:"expire"`
}

type ApplicationAuthApiKeyConfig struct {
	Apikey string `json:"apikey"`
	Expire int64  `json:"expire"`
}

type ApplicationAuthAkSkConfig struct {
	Ak     string `json:"ak"`
	Sk     string `json:"sk"`
	Expire int64  `json:"expire"`
}

type ApplicationAuthJwtConfig struct {
	Iss               string   `json:"iss"`            //	签发机构
	Secret            string   `json:"secret"`         //密钥，当算法是HS时必填
	RsaPublicKey      string   `json:"rsa_public_key"` //RSA公钥，当算法是RS/ES时必填
	Algorithm         string   `json:"algorithm"`
	ClaimsToVerify    []string `json:"claims_to_verify"`
	SignatureIsBase64 bool     `json:"signature_is_base_64"` //signature是否经过base64加密
	Path              string   `json:"path"`                 //获取用户信息字段
}

//type StrategyInfo struct {
//	Name     string              `json:"name"`
//	Stop     bool                `json:"stop"`
//	Priority int                 `json:"priority"`
//	Desc     string              `json:"description"`
//	Filters  map[string][]string `json:"filters"`
//	Limiting StrategyLimiting    `json:"limiting"`
//}
//
//type StrategyLimit struct {
//	Second int `json:"second"`
//	Minute int `json:"minute"`
//	Hour   int `json:"hour"`
//}
//
//type StrategyLimiting struct {
//	Metrics []string      `json:"metrics"`
//	Query   StrategyLimit `json:"query"`
//	Traffic StrategyLimit `json:"traffic"`
//}

// StrategyGrey apinto灰度策略配置
type StrategyGrey struct {
	KeepSession  bool         `json:"keep_session"`
	Nodes        []string     `json:"nodes"`
	Distribution string       `json:"distribution"`
	Percent      int          `json:"percent"`
	Match        []RouterRule `json:"matching"`
}

// StrategyVisit apinto灰度策略配置
type StrategyVisit struct {
	VisitRule       string              `json:"visit_rule"`
	InfluenceSphere map[string][]string `json:"influence_sphere"`
	Continue        bool                `json:"continue"`
}

// RedisOutput Redis配置
type RedisOutput struct {
	OutputConfig
	Addrs    []string `json:"addrs,omitempty"`
	Username string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
}

// InfluxV2Output InfluxV2配置
type InfluxV2Output struct {
	OutputConfig
	Addr   string   `json:"url"`
	Org    string   `json:"org"`
	Bucket string   `json:"bucket"`
	Token  string   `json:"token"`
	Scopes []string `json:"scopes"`
}

type CertConfig struct {
	Name   string `json:"name"`
	Key    string `json:"key"`
	Pem    string `json:"pem"`
	Driver string `json:"driver"`
}

type ExtenderListItem struct {
	Group   string `json:"group"`
	Project string `json:"project"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Id      string `json:"id"`
}

type Render []byte

func (r *Render) UnmarshalJSON(b []byte) error {
	*r = b
	return nil
}

type ExtenderInfo struct {
	Group   string `json:"group"`
	Project string `json:"project"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Render  Render `json:"render"`
}

type PluginTemplateConfig struct {
	Plugins     map[string]*Plugin `json:"plugins,omitempty"` //插件
	Name        string             `json:"name"`
	Driver      string             `json:"driver"`
	Description string             `json:"description"`
}

type PluginTemplateInfo struct {
	Create      string             `json:"create"`
	Description string             `json:"description"`
	Driver      string             `json:"driver"`
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Plugins     map[string]*Plugin `json:"plugins,omitempty"` //插件
	Profession  string             `json:"profession"`
	Update      string             `json:"update"`
}
