package dto

import "time"

type MonitorConfigProxy []byte

func (c *MonitorConfigProxy) MarshalJSON() ([]byte, error) {
	return *c, nil
}

func (c *MonitorConfigProxy) UnmarshalJSON(bytes []byte) error {
	*c = bytes
	return nil
}

func (c *MonitorConfigProxy) String() string {
	return string(*c)
}

type MonitorPartitionInfoProxy struct {
	Name         string             `json:"name"`
	SourceType   string             `json:"source_type"`
	Config       MonitorConfigProxy `json:"config"`
	Env          string             `json:"env"`
	ClusterNames []string           `json:"cluster_names"`
}

// MonCallCountOutput 调用量统计趋势图结构体 分区总览-调用量统计/调用趋势图
type MonCallCountOutput struct {
	Date         []time.Time `json:"date"`
	Status5XX    []int64     `json:"status_5xx"`
	Status4XX    []int64     `json:"status_4xx"`
	ProxyRate    []float64   `json:"proxy_rate"`
	ProxyTotal   []int64     `json:"proxy_total"`
	RequestRate  []float64   `json:"request_rate"`
	RequestTotal []int64     `json:"request_total"`
}

// MonCommonInput 通用请求体 query中读取
type MonCommonInput struct {
	StartTime   int64    `json:"start_time"`
	EndTime     int64    `json:"end_time"`
	PartitionId string   `json:"partition_id"`
	Path        string   `json:"path"`
	Clusters    []string `json:"clusters"`
	Ip          string   `json:"ip"`
	ServiceName string   `json:"service_name"`
	AppId       string   `json:"app_id"`
	ApiId       string   `json:"api_id"`
	Addr        string   `json:"addr"`
	ProxyPath   string   `json:"proxy_path"`
	Services    []string `json:"services"`
	AppIds      []string `json:"app_ids"`
	ApiIds      []string `json:"api_ids"`
}

// MonCommonStatisticsOutput 调用统计
type MonCommonStatisticsOutput struct {
	ApiId       string `json:"api_id"`       //apiID
	ApiName     string `json:"api_name"`     //api名称
	ServiceName string `json:"service_name"` //上游服务名称
	AppName     string `json:"app_name"`     //应用名称
	AppId       string `json:"app_id"`       //应用ID
	Path        string `json:"path"`         //路径
	ProxyPath   string `json:"proxy_path"`   //转发路径
	Ip          string `json:"ip"`           //IP
	Addr        string `json:"addr"`         //目标节点
	IsRed       bool   `json:"is_red"`       //是否标红
	MonCommonData
}

// MonCommonData 通用字段
type MonCommonData struct {
	RequestTotal   int64   `json:"request_total"`   //请求总数
	RequestSuccess int64   `json:"request_success"` //请求成功数
	RequestRate    float64 `json:"request_rate"`    //请求成功率
	ProxyTotal     int64   `json:"proxy_total"`     //转发总数
	ProxySuccess   int64   `json:"proxy_success"`   //转发成功数
	ProxyRate      float64 `json:"proxy_rate"`      //转发成功率
	StatusFail     int64   `json:"status_fail"`     //失败状态数
	AvgResp        float64 `json:"avg_resp"`        //平均响应时间
	MaxResp        int64   `json:"max_resp"`        //最大响应时间
	MinResp        int64   `json:"min_resp"`        //最小响应时间
	AvgTraffic     float64 `json:"avg_traffic"`     //平均流量
	MaxTraffic     float64 `json:"max_traffic"`     //最大流量
	MinTraffic     float64 `json:"min_traffic"`     //最小流量
}

// MonSummaryOutput 请求/转发统计
type MonSummaryOutput struct {
	Total     int `json:"total"`      // 请求总数
	Success   int `json:"success"`    //请求成功数
	Fail      int `json:"fail"`       //请求失败数
	Status4Xx int `json:"status_4xx"` //状态码4xx数
	Status5Xx int `json:"status_5xx"` //状态码5xx数
}

type MonMessageOutput struct {
	Date     []interface{} `json:"date"`
	Request  []interface{} `json:"request"`  //请求报文量
	Response []interface{} `json:"response"` //响应报文量
}

type CircularDate struct {
	Total     int64 `json:"total"`
	Success   int64 `json:"success"`
	Fail      int64 `json:"fail"`
	Status4Xx int64 `json:"status_4xx"`
	Status5Xx int64 `json:"status_5xx"`
}

// MonSummaryInput 总览接口通用输入
type MonSummaryInput struct {
	UUid     string   `json:"uuid"`
	Clusters []string `json:"clusters"`
	Start    int64    `json:"start"`
	End      int64    `json:"end"`
	DataType string   `json:"data_type"`
}
