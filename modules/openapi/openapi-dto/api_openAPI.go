package openapi_dto

type SyncImportData struct {
	Format      string            `json:"format"`
	Content     []byte            `json:"content"`
	ServiceName string            `json:"upstream"`
	Server      *ImportServerInfo `json:"server,omitempty"`
	GroupUUID   string            `json:"group"`
	Label       string            `json:"label"`
	Prefix      string            `json:"prefix"`
}

type ImportServerInfo struct {
	Scheme string             `json:"scheme"`
	Nodes  []*ImportNodesInfo `json:"nodes"`
}

type ImportNodesInfo struct {
	Url    string `json:"url"`
	Weight int    `json:"weight"`
}
