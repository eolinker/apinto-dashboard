/*
 * Copyright (c) 2023. Lorem ipsum dolor sit amet, consectetur adipiscing elit.
 * Morbi non lorem porttitor neque feugiat blandit. Ut vitae ipsum eget quam lacinia accumsan.
 * Etiam sed turpis ac ipsum condimentum fringilla. Maecenas magna.
 * Proin dapibus sapien vel ante. Aliquam erat volutpat. Pellentesque sagittis ligula eget metus.
 * Vestibulum commodo. Ut rhoncus gravida arcu.
 */

package navigation

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
