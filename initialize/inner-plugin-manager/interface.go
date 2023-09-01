package inner_plugin_manager

type IInnerPlugin interface {
	GetDefine() interface{}
	GetIcon(fileName string) ([]byte, error)
	GetReadme(fileName string) ([]byte, error)
	GetResourcesFile(filePath string) ([]byte, error)
}
