package driver_manager

import "sync"

type DriverInfo struct {
	Name   string
	Render string
}

type IDriverManager[T any] interface {
	RegisterDriver(driverName string, driver T)
	GetDriver(driverName string) T
	DelDriver(driverName string)
}

type driverManager[T any] struct {
	mutex   *sync.Mutex
	drivers map[string]T
}

func (d *driverManager[T]) RegisterDriver(driverName string, t T) {
	d.mutex.Lock()
	d.drivers[driverName] = t
	d.mutex.Unlock()
}
func (d *driverManager[T]) GetDriver(driverName string) T {
	return d.drivers[driverName]
}

func (d *driverManager[T]) DelDriver(driverName string) {
	d.mutex.Lock()
	delete(d.drivers, driverName)
	d.mutex.Unlock()
}
func createDriverManager[T any]() *driverManager[T] {
	manager := &driverManager[T]{
		mutex:   new(sync.Mutex),
		drivers: make(map[string]T),
	}
	return manager
}

//type ManagerFactroy func(profession string) IDriverManager

//type IDriverManager interface {
//	registerDiscoveryDriver(driverName string, driver driver.IDiscoveryDriver)
//	registerAuthDriver(driverName string, driver driver.IAuthDriver)
//	GetDiscoveryDriver(driverName string) driver.IDiscoveryDriver
//	GetAuthDriver(driverName string) driver.IAuthDriver
//	DiscoveryList() []*DriverInfo
//	AuthList() []*DriverInfo
//}
//
//type driverManager struct {
//	discoveryDrivers map[string]driver.IDiscoveryDriver
//	authDrivers      map[string]driver.IAuthDriver
//}
//
//func (d *driverManager) GetAuthDriver(driverName string) driver.IAuthDriver {
//	return d.authDrivers[driverName]
//}
//
//func newDriverManager() IDriverManager {
//	manager := &driverManager{
//		discoveryDrivers: make(map[string]driver.IDiscoveryDriver),
//		authDrivers:      make(map[string]driver.IAuthDriver),
//	}
//	return manager
//}
//
//func (d *driverManager) registerDiscoveryDriver(driverName string, driver driver.IDiscoveryDriver) {
//	d.discoveryDrivers[driverName] = driver
//}
//
//func (d *driverManager) registerAuthDriver(driverName string, driver driver.IAuthDriver) {
//	d.authDrivers[driverName] = driver
//}
//
//func (d *driverManager) GetDiscoveryDriver(driverName string) driver.IDiscoveryDriver {
//	return d.discoveryDrivers[driverName]
//}
//
//func (d *driverManager) DiscoveryList() []*DriverInfo {
//	infos := make([]*DriverInfo, 0, len(d.discoveryDrivers))
//	for driverName, dd := range d.discoveryDrivers {
//		info := &DriverInfo{
//			Name:   driverName,
//			Render: dd.Render(),
//		}
//		infos = append(infos, info)
//	}
//	return infos
//}
//
//func (d *driverManager) AuthList() []*DriverInfo {
//	infos := make([]*DriverInfo, 0, len(d.discoveryDrivers))
//	for driverName, dd := range d.authDrivers {
//		info := &DriverInfo{
//			Name:   driverName,
//			Render: dd.Render(),
//		}
//		infos = append(infos, info)
//	}
//	return infos
//}
