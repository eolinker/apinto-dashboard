package apinto

type IAdmin interface {
	NodeAddr() string
	UpdateNodes(nodes []string) string
}

type IClient interface {
	List(profession string) (data []byte, code int, err error)
	Get(profession string, name string) (data []byte, code int, err error)
	Create(profession string, body []byte) (data []byte, code int, err error)
	Delete(profession string, name string) (data []byte, code int, err error)
	Update(profession string, name string, body []byte) (data []byte, code int, err error)
	Patch(profession string, name string, body []byte) (data []byte, code int, err error)
	PatchPath(profession string, name string, path string, body []byte) (data []byte, code int, err error)
}

type IRender interface {
	Render(profession string, driver string) (data []byte, code int, err error)
}
