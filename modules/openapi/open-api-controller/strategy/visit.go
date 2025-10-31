package strategy

import (
	"encoding/json"
	"fmt"
)

type VisitStrategyValidator struct {
}

func (v *VisitStrategyValidator) Name() string {
	return "visit"
}

func (v *VisitStrategyValidator) Validate(data []byte) ([]byte, error) {
	var cfg VisitStrategy
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
	for k, v := range cfg.Config.InfluenceSphere {
		checker, has := GetChecker(k)
		if !has {
			return nil, fmt.Errorf("influence %s not found", k)
		}
		err = checker.Check(v)
		if err != nil {
			return nil, err
		}
	}
	if cfg.Config.VisitRule != "allow" && cfg.Config.VisitRule != "refuse" {
		return nil, fmt.Errorf("invalid visit_rule: %s", cfg.Config.VisitRule)
	}
	return json.Marshal(cfg)
}

type VisitStrategy struct {
	Filters []*Filter           `json:"filters"`
	Config  VisitStrategyConfig `json:"config"`
}

type VisitStrategyConfig struct {
	InfluenceSphere map[string][]string `json:"influence_sphere"`
	VisitRule       string              `json:"visit_rule"`
	Continue        bool                `json:"continue"`
}
