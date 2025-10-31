package pm3

type Cargo struct {
	Value string
	Title string
}

func (c *Cargo) Export() *CargoItem {
	return &CargoItem{
		Value: c.Value,
		Title: c.Title,
	}
}

type CargoItem struct {
	Value string `json:"name"`
	Title string `json:"title"`
}
type CargoStatus int

const (
	None    = iota // 无效
	Offline        // 未上线
	Online         // 已上线
)

type Provider interface {
	Provide(namespaceId int) []Cargo
}
type ProviderStatus interface {
	Status(key string, namespaceId int, cluster string) (CargoStatus, string)
}
type ProviderSupport interface {
	Provider() map[string]Provider
	ProviderStatus
}
