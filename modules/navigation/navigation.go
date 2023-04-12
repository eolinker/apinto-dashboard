package navigation

import (
	navigation_model "github.com/eolinker/apinto-dashboard/modules/navigation/navigation-model"
)

type INavigationService interface {
	Info(uuid string) (*navigation_model.Navigation, bool)
	List() []*navigation_model.Navigation
}

type INavigationDataService interface {
	GetNavigationData() []*navigation_model.Navigation
}
