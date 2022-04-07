package apinto

type IClient interface {
	List(profession string) (*response, error)
	Get(profession string, name string) (*response, error)
	Create(profession string, data []byte) (*response, error)
	Delete(profession string, name string) (*response, error)
	Update(profession string, name string, data []byte) (*response, error)
	Patch(profession string, name string, body []byte) (*response, error)
	Enable(profession string, name string) (*response, error)
}

type IRender interface {
	Render(profession string, driver string) (interface{}, error)
}
