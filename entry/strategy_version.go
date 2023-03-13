package entry

import "time"

type StrategyVersion struct {
	Id          int `json:"id"`
	StrategyId  int `json:"strategy_id"`
	NamespaceId int `json:"namespace_id"`
	StrategyConfigInfo
	Operator   int       `json:"operator"`
	CreateTime time.Time `json:"create_time"`
}

func (v *StrategyVersion) SetVersionId(id int) {
	v.Id = id
}

// StrategyVisitVersionConfig 访问策略config
type StrategyVisitVersionConfig struct {
}

//StrategyFuseVersionConfig 熔断策略config
type StrategyFuseVersionConfig struct {
}

type StrategyConfigInfo struct {
	Priority int                     `json:"priority,omitempty"`
	IsStop   bool                    `json:"is_stop"`
	Type     string                  `json:"type"`
	Filters  []StrategyFiltersConfig `json:"filters,omitempty"`
	StrategyVersionConfig
}

//StrategyVersionConfig 策略config
type StrategyVersionConfig struct {
	Config string `json:"config,omitempty"`
}

type StrategyFiltersConfig struct {
	Name   string   `json:"name,omitempty"`
	Values []string `json:"values,omitempty"`
}

//StrategyTrafficLimitConfig 流量策略配置
type StrategyTrafficLimitConfig struct {
	Metrics  []string             `json:"metrics"` //限流维度
	Query    Limit                `json:"query"`
	Traffic  Limit                `json:"traffic"`
	Response StrategyResponseConf `json:"response"`
}

type Limit struct {
	Second int `json:"second"`
	Minute int `json:"minute"`
	Hour   int `json:"hour"`
}

//StrategyCacheConfig 缓存策略配置
type StrategyCacheConfig struct {
	ValidTime uint64 `json:"valid_time"`
}

//StrategyGreyConfig 灰度策略配置
type StrategyGreyConfig struct {
	KeepSession  bool        `json:"keep_session"`
	Nodes        []string    `json:"nodes"`
	Distribution string      `json:"distribution"`
	Percent      int         `json:"percent"`
	Match        []MatchConf `json:"match"`
}

//StrategyVisitConfig 访问策略配置
type StrategyVisitConfig struct {
	VisitRule       string                  `json:"visit_rule"`
	InfluenceSphere []StrategyFiltersConfig `json:"influence_sphere"`
	Continue        bool                    `json:"continue"`
}

//StrategyFuseConfig 熔断策略配置
type StrategyFuseConfig struct {
	Metric           string               `json:"metric"`         //熔断维度
	FuseCondition    StatusConditionConf  `json:"fuse_condition"` //熔断条件
	FuseTime         FuseTimeConf         `json:"fuse_time"`
	RecoverCondition StatusConditionConf  `json:"recover_condition"` //恢复条件
	Response         StrategyResponseConf `json:"response"`
}

type StatusConditionConf struct {
	StatusCodes []int `json:"status_codes"`
	Count       int   `json:"count"`
}

type FuseTimeConf struct {
	Time    int `json:"time"`
	MaxTime int `json:"max_time"`
}

//StrategyResponseConf 策略返回内容配置
type StrategyResponseConf struct {
	StatusCode  int                      `json:"status_code"`
	ContentType string                   `json:"content_type"`
	Charset     string                   `json:"charset"`
	Header      []StrategyResponseHeader `json:"header,omitempty"`
	Body        string                   `json:"body"`
}

type StrategyResponseHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
