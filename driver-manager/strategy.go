package driver_manager

//type IStrategyDriverManager interface {
//	IDriverManager[driver.IStrategyDriver]
//	List() []*StrategyDriverInfo
//}

//type StrategyDriverInfo struct {
//	Name       string
//	ApintoName string
//}

//type strategyDriver struct {
//	*driverManager[driver.IStrategyDriver]
//}

//func (d *strategyDriver) List() []*StrategyDriverInfo {
//	list := make([]*StrategyDriverInfo, 0)
//	for name, _ := range d.drivers {
//		list = append(list, &StrategyDriverInfo{
//			Name: name,
//		})
//	}
//	return list
//}
//
//func newStrategyDriverManager() IStrategyDriverManager {
//	return &strategyDriver{driverManager: createDriverManager[driver.IStrategyDriver]()}
//}
