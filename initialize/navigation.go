package initialize

import (
	_ "embed"

	"github.com/eolinker/eosc/common/bean"

	navigation_model "github.com/eolinker/apinto-dashboard/modules/navigation/navigation-model"

	"gopkg.in/yaml.v3"
)

var (
	//go:embed navigation.yml
	navigationContent []byte
)

func InitNavigation() error {
	// 初始化导航
	navigations := make([]*navigation_model.Navigation, 0)
	err := yaml.Unmarshal(navigationContent, &navigations)
	if err != nil {
		return err
	}
	service := newNavigationDataService(navigations)
	bean.Injection(&service)
	return nil
}

type navigationDataService struct {
	navigations []*navigation_model.Navigation
}

func newNavigationDataService(navigations []*navigation_model.Navigation) *navigationDataService {
	return &navigationDataService{navigations: navigations}
}

func (n *navigationDataService) GetNavigationData() []*navigation_model.Navigation {
	return n.navigations
}
