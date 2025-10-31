package pm3

type PModule struct {
	Navigation string
	Name       string
	Cname      string
	Router     string
}
type PAccess struct {
	Name   string
	Cname  string
	Module string
	Depend []string
}
type PFrontend map[string]interface{}

type PluginConfig map[string][]ExtendParams
type ExtendParams struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
