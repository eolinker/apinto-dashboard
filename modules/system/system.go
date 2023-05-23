package system

type ISystemConfigService[T any] interface {
	Get() (*T, error)
	Set(t *T) error
}
type ISystemConfigServiceSimple interface {
	Get() ([]byte, error)
	Set(v []byte) error
}
type ISystemConfigDataService interface {
	Get(key string) ([]byte, error)
	Set(key string, v []byte) error
}
