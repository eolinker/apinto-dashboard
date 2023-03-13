package entry

type QuoteKindType string

var (
	QuoteKindTypeService      QuoteKindType = "service"
	QuoteKindTypeDiscovery    QuoteKindType = "discovery"
	QuoteKindTypeAPI          QuoteKindType = "api"
	QuoteKindTypeWarnStrategy QuoteKindType = "warn_strategy"
)

type QuoteTargetKindType string

var (
	QuoteTargetKindTypeVariable      QuoteTargetKindType = "variable"
	QuoteTargetKindTypeDiscovery     QuoteTargetKindType = "discovery"
	QuoteTargetKindTypeService       QuoteTargetKindType = "service"
	QuoteTargetKindTypeNoticeChannel QuoteTargetKindType = "notice_channel"
)

// Quote 引用表
type Quote struct {
	ID         int                 `gorm:"type:int(11);size:11;not null;auto_increment;primary_key;column:id;comment:主键ID"`
	Kind       QuoteKindType       `gorm:"size:20;not null;column:kind;index:kind;dbUniqueIndex:unique;comment:类型"`
	Source     int                 `gorm:"type:int(11);size:11;not null;column:source;dbUniqueIndex:unique;comment:根据kind区分是哪个表的ID,被引用的id"`
	Target     int                 `gorm:"type:int(11);size:11;not null;column:target;comment:引用的id"`
	TargetKind QuoteTargetKindType `gorm:"size:20;not null;index:target_kind;column:target_kind;comment:引用的类型"`
}

func (s *Quote) IdValue() int {
	return s.ID
}

func (*Quote) TableName() string {
	return "quote"
}
