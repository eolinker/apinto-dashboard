package initialize

import (
	_ "embed"
)

var (
	//go:embed navigation.yml
	navigationContent []byte
)

type NavigationCfg struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
	ICON string `yaml:"icon"`
}

func InitNavigation() error {
	// 初始化导航
	//var service navigation.INavigationService
	//bean.Autowired(&service)
	//ctx := context.Background()
	//navigations := make([]*NavigationCfg, 0)
	//err := yaml.Unmarshal(navigationContent, &navigations)
	//if err != nil {
	//	return err
	//}
	//for _, nav := range navigations {
	//	_, err = service.GetIDByUUID(ctx, nav.ID)
	//	if err != nil {
	//		if err == gorm.ErrRecordNotFound {
	//			err = service.Add(ctx, nav.ID, nav.Name, nav.ICON, "css")
	//			if err != nil {
	//				return err
	//			}
	//		}
	//		return err
	//	}
	//}
	return nil
}
