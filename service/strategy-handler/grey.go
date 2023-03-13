package strategy_handler

import (
	"errors"
	"fmt"
	v1 "github.com/eolinker/apinto-dashboard/client/v1"
	"github.com/eolinker/apinto-dashboard/common"
	"github.com/eolinker/apinto-dashboard/dto"
	"github.com/eolinker/apinto-dashboard/entry"
	"github.com/eolinker/apinto-dashboard/enum"
	"github.com/eolinker/apinto-dashboard/service"
	"net/textproto"
	"strings"
)

type greyHandler struct {
	service.FormatHandler[entry.StrategyGreyConfig]
	apintoDriverName string
}

func (t *greyHandler) GetListLabel(conf *entry.StrategyGreyConfig) string {
	switch conf.Distribution {
	case enum.GreyDistributionPercent:
		return "按百分比"
	case enum.GreyDistributionMatch:
		return "按规则"
	default:
		return ""
	}
}

func (t *greyHandler) GetType() string {
	return enum.StrategyGrey
}

func (t *greyHandler) GetConfName() string {
	return enum.StrategyGreyApintoConfName
}

// GetBatchSettingName 获取往apinto发送批量操作策略时 url所需要的路径名 /setting/xxx
func (t *greyHandler) GetBatchSettingName() string {
	return enum.StrategyGreyBatchName
}

func (t *greyHandler) CheckInput(input *dto.StrategyInfoInput[entry.StrategyGreyConfig]) error {
	input.Uuid = strings.TrimSpace(input.Uuid)
	if input.Uuid != "" {
		err := common.IsMatchString(common.UUIDExp, input.Uuid)
		if err != nil {
			return err
		}
	}

	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		return errors.New("Name can't be null. ")
	}
	if input.Priority < 0 {
		input.Priority = 0
	}

	if input.Config == nil {
		return errors.New("config can't be null. ")
	}

	//校验百分比
	if input.Config.Percent > 10000 || input.Config.Percent < 0 {
		return fmt.Errorf("percent %d is illegal. ", input.Config.Percent)
	}

	//校验灰度节点
	nodes := input.Config.Nodes
	if len(nodes) == 0 {
		return errors.New("nodes can't be null. ")
	}

	for _, node := range input.Config.Nodes {
		//若节点地址不符合域名:port或则ip:port
		if !common.IsMatchDomainPort(node) && !common.IsMatchIpPort(node) {
			return fmt.Errorf("node addr %s is illegal. ", node)
		}
	}

	//校验分配方式
	switch input.Config.Distribution {
	case enum.GreyDistributionPercent, enum.GreyDistributionMatch:
	default:
		return fmt.Errorf("distribution %s is illegal. ", input.Config.Distribution)
	}

	//校验规则
	for i, m := range input.Config.Match {
		input.Config.Match[i].Key = strings.TrimSpace(m.Key)
		if input.Config.Match[i].Key == "" {
			return errors.New("Match.Key can't be nil. ")
		}
		switch m.Position {
		case enum.MatchPositionHeader, enum.MatchPositionQuery, enum.MatchPositionCookie:
		default:
			return fmt.Errorf("position %s is illegal. ", m.Position)
		}

		switch m.MatchType {
		case enum.MatchTypeEqual, enum.MatchTypePrefix, enum.MatchTypeSuffix, enum.MatchTypeSubstr, enum.MatchTypeUneuqal, enum.MatchTypeRegexp, enum.MatchTypeRegexpG:
			input.Config.Match[i].Pattern = strings.TrimSpace(m.Pattern)
			if input.Config.Match[i].Pattern == "" {
				return fmt.Errorf("Match.Pattern can't be nil when MatchType is %s. ", m.MatchType)
			}
		case enum.MatchTypeNull, enum.MatchTypeExist, enum.MatchTypeUnexist, enum.MatchTypeAny:
		default:
			return fmt.Errorf("match_type %s is illegal. ", m.MatchType)
		}
	}

	//校验筛选条件
	return checkFilters(input.Filters)
}

func (t *greyHandler) ToApintoConfig(conf entry.StrategyGreyConfig) interface{} {
	rules := make([]v1.RouterRule, 0, len(conf.Match))

	for _, m := range conf.Match {
		rule := v1.RouterRule{
			Type:  m.Position,
			Name:  m.Key,
			Value: "",
		}

		if m.Position == enum.MatchPositionHeader {
			rule.Name = textproto.CanonicalMIMEHeaderKey(rule.Name)
		}

		switch m.MatchType {
		case enum.MatchTypeEqual:
			rule.Value = m.Pattern
		case enum.MatchTypePrefix:
			rule.Value = fmt.Sprintf("%s*", m.Pattern)
		case enum.MatchTypeSuffix:
			rule.Value = fmt.Sprintf("*%s", m.Pattern)
		case enum.MatchTypeSubstr:
			rule.Value = fmt.Sprintf("*%s*", m.Pattern)
		case enum.MatchTypeUneuqal:
			rule.Value = fmt.Sprintf("!=%s", m.Pattern)
		case enum.MatchTypeNull:
			rule.Value = "$"
		case enum.MatchTypeExist:
			rule.Value = "**"
		case enum.MatchTypeUnexist:
			rule.Value = "!"
		case enum.MatchTypeRegexp:
			rule.Value = fmt.Sprintf("~=%s", m.Pattern)
		case enum.MatchTypeRegexpG:
			rule.Value = fmt.Sprintf("~*=%s", m.Pattern)
		case enum.MatchTypeAny:
			rule.Value = "*"
		}

		rules = append(rules, rule)
	}

	apintoConf := &v1.StrategyGrey{
		KeepSession:  conf.KeepSession,
		Nodes:        conf.Nodes,
		Distribution: conf.Distribution,
		Percent:      conf.Percent,
		Match:        rules,
	}

	return apintoConf
}

func NewStrategyGreyHandler(apintoDriverName string) service.IStrategyHandler[entry.StrategyGreyConfig, entry.StrategyGreyConfig] {
	return &greyHandler{
		apintoDriverName: apintoDriverName,
	}
}
