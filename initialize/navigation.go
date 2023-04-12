package initialize

import (
	_ "embed"

	"github.com/eolinker/apinto-dashboard/modules/navigation"

	"github.com/eolinker/eosc/common/bean"

	navigation_model "github.com/eolinker/apinto-dashboard/modules/navigation/navigation-model"

	"gopkg.in/yaml.v3"
)

var (
	//go:embed navigation.yml
	navigationContent []byte
)

func init() {
	// 初始化导航
	navigations := make([]*navigation_model.Navigation, 0)
	err := yaml.Unmarshal(navigationContent, &navigations)
	if err != nil {
		panic(err)
	}
	service := newNavigationDataService(navigations)
	bean.Injection(&service)
}

type navigationDataService struct {
	navigations []*navigation_model.Navigation
}

func newNavigationDataService(navigations []*navigation_model.Navigation) navigation.INavigationDataService {
	return &navigationDataService{navigations: navigations}
}

func (n *navigationDataService) GetNavigationData() []*navigation_model.Navigation {
	return n.navigations
}
