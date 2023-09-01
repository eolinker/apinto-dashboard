package dynamic_controller

import (
	"context"
	"fmt"

	apinto_module "github.com/eolinker/apinto-dashboard/module"

	"github.com/eolinker/apinto-dashboard/modules/dynamic"
	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/log"
)

type skillProvider struct {
	profession     string
	skill          string
	dynamicService dynamic.IDynamicService
}

func newSkillProvider(profession string, skill string) *skillProvider {
	p := &skillProvider{profession: profession, skill: skill}
	bean.Autowired(&p.dynamicService)
	return p
}

func (s *skillProvider) Provide(namespaceID int) []apinto_module.Cargo {
	list, err := s.dynamicService.GetBySkill(context.Background(), namespaceID, s.skill)
	if err != nil {
		log.Error("get skill error: ", err)
		return nil
	}
	result := make([]apinto_module.Cargo, 0, len(list))
	for _, l := range list {
		result = append(result, apinto_module.Cargo{
			Value: l.ID,
			Title: fmt.Sprintf("%s[%s]", l.Title, l.Driver),
		})
	}
	return result
}
