package dto

type WebhookInput struct {
	Uuid          string            `json:"uuid"`
	Title         string            `json:"title"`
	Desc          string            `json:"desc"`
	Url           string            `json:"url"`
	Method        string            `json:"method"`
	ContentType   string            `json:"content_type"` //
	NoticeType    string            `json:"notice_type"`  //单次：single  多次：many
	UserSeparator string            `json:"user_separator"`
	Header        map[string]string `json:"header"`
	Template      string            `json:"template"`
}

type WebhooksOutput struct {
	Uuid        string `json:"uuid"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	Method      string `json:"method"`
	ContentType string `json:"content_type"`
	Operator    string `json:"operator"`
	UpdateTime  string `json:"update_time"`
	CreateTime  string `json:"create_time"`
	IsDelete    bool   `json:"is_delete"`
}

type WebhookOutput struct {
	Uuid          string            `json:"uuid"`
	Title         string            `json:"title"`
	Desc          string            `json:"desc"`
	Url           string            `json:"url"`
	Method        string            `json:"method"`
	ContentType   string            `json:"content_type"` //
	NoticeType    string            `json:"notice_type"`  //单次：single  多次：many
	UserSeparator string            `json:"user_separator"`
	Header        map[string]string `json:"header"`
	Template      string            `json:"template"`
}

type EmailInput struct {
	Uuid     string `json:"uuid"`
	SmtpUrl  string `json:"smtp_url"`
	SmtpPort int    `json:"smtp_port"`
	Protocol string `json:"protocol"`
	Email    string `json:"email"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

type EmailOutput struct {
	Uuid     string `json:"uuid"`
	SmtpUrl  string `json:"smtp_url"`
	SmtpPort int    `json:"smtp_port"`
	Protocol string `json:"protocol"`
	Email    string `json:"email"`
	Account  string `json:"account"`
	Password string `json:"password"`
}

type WarnHistory struct {
	StrategyTitle string `json:"strategy_title"` //策略名称
	WarnTarget    string `json:"warn_target"`    //告警目标
	WarnContent   string `json:"warn_content"`   //告警内容
	CreateTime    string `json:"create_time"`
	Status        int    `json:"status"`
	ErrMsg        string `json:"err_msg"`
}

type NoticeChannel struct {
	Uuid  string `json:"uuid"`
	Title string `json:"title"`
	Type  int    `json:"type"` //1邮箱 2webhook
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
