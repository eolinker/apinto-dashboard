package strategy

import (
	"encoding/json"
	"fmt"
)

type LimitStrategyValidator struct {
}

func (l *LimitStrategyValidator) Name() string {
	return "limiting"
}

func (l *LimitStrategyValidator) Validate(data []byte) ([]byte, error) {
	var cfg LimitStrategy
	err := json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}
	if len(cfg.Filters) < 1 {
		return nil, fmt.Errorf("filters not found")
	}
	for _, f := range cfg.Filters {
		checker, has := GetChecker(f.Name)
		if !has {
			return nil, fmt.Errorf("filter %s not found", f.Name)
		}
		err = checker.Check(f.Values)
		if err != nil {
			return nil, err
		}
	}
	if len(cfg.Config.Metrics) < 1 {
		return nil, fmt.Errorf("metrics not found")
	}
	for _, metric := range cfg.Config.Metrics {
		switch metric {
		case "{api}", "{ip}", "{service}":
		default:
			return nil, fmt.Errorf("invalid metric: %s", metric)
		}
	}
	return json.Marshal(cfg)
}

type LimitStrategy struct {
	Filters []*Filter           `json:"filters"`
	Config  LimitStrategyConfig `json:"config"`
}

type LimitStrategyConfig struct {
	Metrics  []string        `json:"metrics"`
	Query    LimitQuery      `json:"query"`
	Response RewriteResponse `json:"response"`
	Traffic  LimitTraffic    `json:"traffic"`
}

type LimitQuery struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`
}

type LimitTraffic struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
	Second int `json:"second"`
}
