package inner_plugin_manager

var _ IInnerPlugin = (*embedPlugin)(nil)

type embedPlugin struct {
	//TODO 存FS?
}

func (e embedPlugin) GetDefine() interface{} {
	//TODO implement me
	panic("implement me")
}

func (e embedPlugin) GetIcon(fileName string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (e embedPlugin) GetReadme(fileName string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (e embedPlugin) GetResourcesFile(filePath string) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}
