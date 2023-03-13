package extend

type ApiConfig []ApiModule
type ApiModule struct {
	ModuleName string        `json:"module"`
	Describe   string        `json:"describe"`
	Functions  []ApiFunction `json:"functions"`
}

type ApiFunction struct {
	Function string `json:"function" yaml:"function"`
	Method   string `json:"method"` //后端接口方法
	Path     string `json:"path"`   // 后端接口path

}
