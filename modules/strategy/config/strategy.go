package config

const (
	StrategyTraffic = "traffic"
	StrategyCache   = "cache"
	StrategyGrey    = "grey"
	StrategyFuse    = "fuse"
	StrategyVisit   = "visit"

	StrategyTrafficApintoConfName = "limiting"
	StrategyCacheApintoConfName   = "valid_time"
	StrategyGreyApintoConfName    = "grey"
	StrategyFuseApintoConfName    = "fuse"
	StrategyVisitApintoConfName   = "visit"

	StrategyTrafficRuntimeKind = "strategy_traffic"
	StrategyCacheRuntimeKind   = "strategy_cache"
	StrategyGreyRuntimeKind    = "strategy_grey"
	StrategyVisitRuntimeKind   = "strategy_visit"
	StrategyFuseRuntimeKind    = "strategy_fuse"

	StrategyTrafficBatchName = "strategies-limiting"
	StrategyCacheBatchName   = "strategies-cache"
	StrategyGreyBatchName    = "strategies-grey"
	StrategyFuseBatchName    = "strategies-fuse"
	StrategyVisitBatchName   = "strategies-visit"

	FilterApplication = "application"
	FilterApi         = "api"
	FilterPath        = "path"
	FilterService     = "service"
	FilterMethod      = "method"
	FilterIP          = "ip"
	FilterAppKey      = "appkey"

	FilterTypeRemote  = "remote"
	FilterTypePattern = "pattern"
	FilterTypeStatic  = "static"

	FilterValuesALL = "ALL"
	ApiPathRegexp   = `^\*?[\w-/]+\*?$`

	MetricsAPP      = "{application}"
	MetricsAPI      = "{api}"
	MetricsService  = "{service}"
	MetricsStrategy = "{strategy}"
	MetricsIP       = "{ip}"

	//灰度策略分配方式
	GreyDistributionPercent = "percent"
	GreyDistributionMatch   = "match"

	//访问策略访问规则
	VisitRuleAllow  = "allow"
	VisitRuleRefuse = "refuse"

	//熔断策略字符集
	CharsetUTF8  = "UTF-8"
	CharsetGBK   = "GBK"
	CharsetASCII = "ASCII"
)

func GetStrategyCharsetList() []string {
	return []string{CharsetUTF8, CharsetGBK, CharsetASCII}
}
