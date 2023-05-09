package warn_dto

type WarnHistory struct {
	StrategyTitle string `json:"strategy_title"` //策略名称
	WarnTarget    string `json:"warn_target"`    //告警目标
	WarnContent   string `json:"warn_content"`   //告警内容
	CreateTime    string `json:"create_time"`
	Status        int    `json:"status"`
	ErrMsg        string `json:"err_msg"`
}

type WarnStrategyList struct {
	Uuid          string `json:"uuid"`
	StrategyTitle string `json:"strategy_title"`
	WarnDimension string `json:"warn_dimension"`
	WarnTarget    string `json:"warn_target"`
	WarnRule      string `json:"warn_rule"`
	WarnFrequency string `json:"warn_frequency"`
	IsEnable      bool   `json:"is_enable"`
	Operator      string `json:"operator"`
	UpdateTime    string `json:"update_time"`
	CreateTime    string `json:"create_time"`
}

type WarnStrategyRule struct {
	ChannelUuids []string                     `json:"channel_uuids"`
	Condition    []*WarnStrategyRuleCondition `json:"condition"`
}

type WarnStrategyRuleCondition struct {
	Compare string  `json:"compare"`
	Unit    string  `json:"unit"`
	Value   float64 `json:"value"`
}

type WarnStrategyTarget struct {
	Rule   string   `json:"rule"`
	Values []string `json:"values"`
}
type WarnStrategy struct {
	Uuid       string              `json:"uuid"`
	Title      string              `json:"title"`
	Desc       string              `json:"desc"`
	IsEnable   bool                `json:"is_enable"`
	Dimension  string              `json:"dimension"`
	Target     *WarnStrategyTarget `json:"target"`
	Quota      string              `json:"quota"`
	Every      int                 `json:"every"`
	Rule       []*WarnStrategyRule `json:"rule"`
	Continuity int                 `json:"continuity"`
	HourMax    int                 `json:"hour_max"`
	Users      []int               `json:"users"`
}

type WarnStrategyInput struct {
	PartitionId string              `json:"partition_id"`
	Uuid        string              `json:"uuid"`
	Title       string              `json:"title"`
	Desc        string              `json:"desc"`
	IsEnable    bool                `json:"is_enable"`
	Dimension   string              `json:"dimension"`
	Target      WarnStrategyTarget  `json:"target"`
	Quota       string              `json:"quota"`
	Every       int                 `json:"every"`
	Rule        []*WarnStrategyRule `json:"rule"`
	Continuity  int                 `json:"continuity"`
	HourMax     int                 `json:"hour_max"`
	Users       []int               `json:"users"`
}
